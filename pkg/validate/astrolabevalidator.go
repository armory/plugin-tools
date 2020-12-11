package validate

import (
	"fmt"
	"log"
)

type AstrolabeCompatibilityValidator struct{}

func (v *AstrolabeCompatibilityValidator) IsPluginCompatibleWithPlatform(platformVersion string, pluginId string, pluginVersion string, repos []PluginRepository) (Verdict, error) {
	var metadata *CompatibilityMetadata
	for _, e := range repos {
		m, err := e.getCompatibilityMetadata(pluginId)
		if err != nil {
			log.Printf("error when geting compatibility metadata: %s", err)
		} else {
			metadata = m
			break
		}
	}

	if metadata == nil {
		return Unknown, fmt.Errorf("compatibility metadata not found in repositories")
	}

	tests, err := metadata.getCompatibilityTests(pluginVersion)
	if err != nil {
		return Unknown, err
	}

	for _, e := range tests {
		if e.containsPlatformVersion(platformVersion) {
			if e.Status == "success" {
				return Compatible, nil
			} else {
				return NotCompatible, nil
			}
		}
	}

	return Unknown, fmt.Errorf("could not find compatibility tests with platform: %s for plugin %s@%s", platformVersion, pluginId, pluginVersion)
}

func (v *AstrolabeCompatibilityValidator) IsPluginCompatibleWithService(serviceName string, serviceVersion string, pluginId string, pluginVersion string, repos []PluginRepository) (Verdict, error) {
	var metadata *CompatibilityMetadata
	for _, e := range repos {
		m, err := e.getCompatibilityMetadata(pluginId)
		if err != nil {
			log.Printf("error message: %s", err)
		} else {
			metadata = m
			break
		}
	}

	if metadata == nil {
		return Unknown, fmt.Errorf("compatibility metadata not found in repositories")
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
