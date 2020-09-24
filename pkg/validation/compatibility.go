package validation

import (
	"encoding/json"
	"fmt"
	"github.com/blang/semver"
	"io/ioutil"
	"log"
	"net/http"
)

type repository struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type pluginMetadata struct {
	PluginId    string    `json:"id"`
	Description string    `json:"description"`
	Provider    string    `json:"provider"`
	Releases    []release `json:"releases"`
}

type release struct {
	Version       string    `json:"version"`
	Date          string    `json:"date"`
	Requires      string    `json:"requires"`
	Compatibility *platform `json:"compatibility"`
	Hash          string    `json:"sha512sum"`
	State         string    `json:"state"`
	Url           string    `json:"url"`
}
type platform struct {
	SpinnakerVersions []string `json:"spinnaker"`
}

type plugin struct {
	Id      string
	Version string
}

type CompatibilityResult struct {
	PluginId      string
	PluginVersion string
	IsCompatible  bool
	Reason        string
}

func ArePluginsCompatible(spinnakerVersion string, plugins []plugin, repo []string) []CompatibilityResult {
	var result []CompatibilityResult
	spinVersion, err := semver.Make(spinnakerVersion)
	if err != nil {
		log.Println(err)
		return result
	}
	pluginsMetadata := getPluginMetadata(repo)
	for _, v := range plugins {
		comp, err := getCompatibilityConstraint(v.Id, v.Version, pluginsMetadata)
		if err != nil {
			log.Println(err)
			result = append(result, CompatibilityResult{v.Id, v.Version, true, err.Error()})
			continue
		}
		if comp == nil {
			message := fmt.Sprintf("Plugin %s@%s does not contain compatibility constrain", v.Id, v.Version)
			log.Printf(message)
			result = append(result, CompatibilityResult{v.Id, v.Version, true, message})
			continue
		}
		isCompatible := false
		compatibleVersion := ""
		for _, s := range comp.SpinnakerVersions {
			compSpinVersion, err := semver.Make(s)
			if err != nil {
				log.Printf("%s in compatibility metadata is invalid", compSpinVersion)
				continue
			}
			if spinVersion.Equals(compSpinVersion) || spinVersion.Minor == compSpinVersion.Minor {
				isCompatible = true
				compatibleVersion = s
				break
			}
		}
		if isCompatible {
			result = append(result, CompatibilityResult{v.Id, v.Version, true, fmt.Sprintf("Plugin %s@%s compatible with version %s", v.Id, v.Version, compatibleVersion)})
		} else {
			result = append(result, CompatibilityResult{v.Id, v.Version, false, fmt.Sprintf("No compatible Spinnaker versions found for Plugin %s@%s", v.Id, v.Version)})
		}
	}
	return result
}

func getPluginMetadata(repositories []string) []pluginMetadata {
	var allPluginsMetadata []pluginMetadata
	for _, s := range repositories {
		body := getExternalResource(s)

		//get list of repositories
		var pf4jRepos []repository
		jsonErr := json.Unmarshal(body, &pf4jRepos)
		if jsonErr != nil {
			log.Println(jsonErr)
			return allPluginsMetadata
		}

		for _, v := range pf4jRepos {
			b := getExternalResource(v.Url)
			var metadata []pluginMetadata
			jsonErr := json.Unmarshal(b, &metadata)
			if jsonErr != nil {
				log.Println(jsonErr)
				continue
			}
			allPluginsMetadata = append(allPluginsMetadata, metadata...)
		}
	}
	return allPluginsMetadata
}

func getExternalResource(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	return body
}

func getCompatibilityConstraint(pluginId string, pluginVersion string, metadata []pluginMetadata) (*platform, error) {
	var releases []release
	for _, v := range metadata {
		if v.PluginId == pluginId {
			releases = v.Releases
			break
		}
	}
	if len(releases) == 0 {
		return &platform{}, fmt.Errorf("No releases found for %s", pluginId)
	}

	for _, v := range releases {
		if v.Version == pluginVersion {
			return v.Compatibility, nil
		}
	}
	return &platform{}, fmt.Errorf("Could not find version %s for plugin %s", pluginVersion, pluginId)
}
