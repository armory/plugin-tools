package validate

import (
	"fmt"
)

type AstrolabeCompatibilityValidator struct{}

func (v *AstrolabeCompatibilityValidator) IsPluginCompatibleWithPlatform(platformVersion string, pluginId string, pluginVersion string, repos []PluginRepository) (Verdict, error) {
	var metadata *CompatibilityMetadata
	for _, e := range repos {
		m, err := e.getCompatibilityMetadata(pluginId)
		if err == nil {
			metadata = m
			break
		}
	}

	if metadata == nil {
		return Unknown, fmt.Errorf("compatibility metadata not found for plugin %s in repositories", pluginId)
	}

	tests, err := metadata.getCompatibilityTests(pluginVersion)
	if err != nil {
		return Unknown, err
	}

	platformFound := false //this flag indicates at least one test was found for given platform
	isUnknown := false
	isCompatible := true
	for _, e := range tests {
		if e.containsPlatformVersion(platformVersion) {
			platformFound = true
			if e.Status == "failure" {
				isCompatible = false
				break
			}
			if e.Status == "unknown" {
				isUnknown = true
				break
			}
		}
	}

	if isUnknown || !platformFound {
		return Unknown, fmt.Errorf("could not find compatibility tests with platform: %s for plugin %s@%s", platformVersion, pluginId, pluginVersion)
	}

	if isCompatible {
		return Compatible, nil
	} else {
		return NotCompatible, nil
	}
}

func (v *AstrolabeCompatibilityValidator) IsPluginCompatibleWithService(serviceName string, serviceVersion string, pluginId string, pluginVersion string, repos []PluginRepository) (Verdict, error) {
	var metadata *CompatibilityMetadata
	for _, e := range repos {
		m, err := e.getCompatibilityMetadata(pluginId)
		if err == nil {
			metadata = m
			break
		}
	}

	if metadata == nil {
		return Unknown, fmt.Errorf("compatibility metadata not found for plugin %s in repositories", pluginId)
	}

	tests, err := metadata.getCompatibilityTests(pluginVersion)
	if err != nil {
		return Unknown, err
	}

	for _, e := range tests {
		if e.ServiceName == serviceName && e.ServiceVersion == serviceVersion {
			if e.Status == "success" {
				return Compatible, nil
			} else {
				return NotCompatible, nil
			}
		}
	}

	return Unknown, fmt.Errorf("could not find compatibility tests with service: %s version: %s for plugin %s@%s", serviceName, serviceVersion, pluginId, pluginVersion)
}
