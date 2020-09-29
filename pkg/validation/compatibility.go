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

func ResolvePluginCompatibility(spinnakerVersion string, plugins []plugin, repo []string) ([]CompatibilityResult, error) {
	var result []CompatibilityResult
	spinVersion, err := semver.Make(spinnakerVersion)
	if err != nil {
		return nil, err
	}
	pluginsMetadata, err := getPluginMetadata(repo)
	if err != nil {
		return nil, err
	}
	for _, v := range plugins {
		comp, err := getCompatibilityConstraint(v.Id, v.Version, pluginsMetadata)
		if err != nil {
			log.Println(err)
			result = append(result, CompatibilityResult{v.Id, v.Version, true, err.Error()})
			continue
		}
		if comp == nil {
			message := fmt.Sprintf("Plugin %s@%s does not contain compatibility constraint", v.Id, v.Version)
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
			if spinVersion.Equals(compSpinVersion) || spinVersion.Minor == compSpinVersion.Minor && spinVersion.Major == compSpinVersion.Major {
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
	return result, nil
}

func getPluginMetadata(repositories []string) ([]pluginMetadata, error) {
	var allPluginsMetadata []pluginMetadata
	for _, s := range repositories {
		body, err := getExternalResource(s)
		if err != nil {
			return nil, err
		}

		//get list of repositories
		var pf4jRepos []repository
		jsonErr := json.Unmarshal(body, &pf4jRepos)
		if jsonErr != nil {
			return nil, err
		}

		for _, v := range pf4jRepos {
			b, err := getExternalResource(v.Url)
			if err != nil {
				return nil, err
			}
			var metadata []pluginMetadata
			jsonErr := json.Unmarshal(b, &metadata)
			if jsonErr != nil {
				return nil, err
			}
			allPluginsMetadata = append(allPluginsMetadata, metadata...)
		}
	}
	return allPluginsMetadata, nil
}

func getExternalResource(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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
