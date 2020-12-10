package validate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type PluginCompatibilityValidator interface {
	IsPluginCompatibleWithPlatform(platformVersion string, pluginId string, pluginVersion string, repos []PluginRepository) (bool, error)
	IsPluginCompatibleWithService(serviceName string, serviceVersion string, pluginId string, pluginVersion string, repos []PluginRepository) (bool, error)
}

func NewAstrolabeValidator() PluginCompatibilityValidator {
	return &AstrolabeCompatibilityValidator{}
}

type Plugin struct {
	id      string
	enabled bool
	version string
}

type PluginRepository struct {
	id  string
	url string
}

func (r *PluginRepository) getCompatibilityMetadata(pluginId string) (*CompatibilityMetadata, error) {
	comUrl := strings.ReplaceAll(r.url, "repositories.json", fmt.Sprintf("compatibility/%s.json", pluginId))
	resp, err := http.Get(comUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("compatibility metadata not found in %s", comUrl)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	metadata := &CompatibilityMetadata{}
	jsonErr := json.Unmarshal(body, &metadata)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return metadata, nil
}

type CompatibilityMetadata struct {
	Id       string                 `json:"id"`
	Releases []CompatibilityRelease `json:"releases"`
}

func (m *CompatibilityMetadata) getCompatibilityTests(pluginVersion string) ([]CompatibilityTest, error) {
	for i, e := range m.Releases {
		if e.Version == pluginVersion {
			return m.Releases[i].Tests, nil
		}
	}
	return nil, fmt.Errorf("release %s not found in compatibility metadata", pluginVersion)
}

type CompatibilityRelease struct {
	Version string              `json:"version"`
	Tests   []CompatibilityTest `json:"compatibility"`
}

type CompatibilityTest struct {
	Status         string   `json:"outcome"`
	Platforms      []string `json:"platformVersions"`
	ServiceName    string   `json:"service"`
	ServiceVersion string   `json:"version"`
}

func (t *CompatibilityTest) containsPlatformVersion(platformVersion string) bool {
	for _, e := range t.Platforms {
		if e == platformVersion {
			return true
		}
	}
	return false
}
