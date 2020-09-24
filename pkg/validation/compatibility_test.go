package validation

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func check(t *testing.T, e error) {
	if e != nil {
		t.Errorf("error when testing: %s", e)
	}
}

func getArrayByteFromPath(location string, t *testing.T) []byte {
	wd, _ := os.Getwd()
	path := fmt.Sprintf("%s/../../%s", wd, location)
	f, err := os.Open(path)
	check(t, err)
	s, err := ioutil.ReadAll(f)
	check(t, err)
	return s
}

func TestCompatiblePluginWithExactVersion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock repository
	httpmock.RegisterResponder("GET", "http://mock.repo/repository",
		httpmock.NewStringResponder(200, `[{"id": "test-repo", "url": "http://mock.repo/test-plugins.json"}]`))

	s := getArrayByteFromPath("test/metadata/valid-metadata.json", t)
	// mock metadata
	httpmock.RegisterResponder("GET", "http://mock.repo/test-plugins.json",
		httpmock.NewBytesResponder(200, s))

	repos := []string{
		"http://mock.repo/repository",
	}
	plugins := []plugin{
		{"Test.pluginA", "1.1.0"},
	}
	expectedResult := []CompatibilityResult{
		{"Test.pluginA", "1.1.0", true, "Plugin Test.pluginA@1.1.0 compatible with version 1.20.0"},
	}
	if got := ArePluginsCompatible("1.20.0", plugins, repos); !reflect.DeepEqual(got, expectedResult) {
		t.Errorf("method() = %v, want %v", got, expectedResult)
	}
}

func TestCompatiblePluginWithMinorVersion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock repository
	httpmock.RegisterResponder("GET", "http://mock.repo/repository",
		httpmock.NewStringResponder(200, `[{"id": "test-repo", "url": "http://mock.repo/test-plugins.json"}]`))

	s := getArrayByteFromPath("test/metadata/valid-metadata.json", t)
	// mock metadata
	httpmock.RegisterResponder("GET", "http://mock.repo/test-plugins.json",
		httpmock.NewBytesResponder(200, s))

	repos := []string{
		"http://mock.repo/repository",
	}
	plugins := []plugin{
		{"Test.pluginA", "1.1.0"},
	}
	expectedResult := []CompatibilityResult{
		{"Test.pluginA", "1.1.0", true, "Plugin Test.pluginA@1.1.0 compatible with version 1.20.0"},
	}
	if got := ArePluginsCompatible("1.20.7", plugins, repos); !reflect.DeepEqual(got, expectedResult) {
		t.Errorf("method() = %v, want %v", got, expectedResult)
	}
}

func TestIncompatiblePluginWithoutConstrainForNewerSpinnakerVersion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock repository
	httpmock.RegisterResponder("GET", "http://mock.repo/repository",
		httpmock.NewStringResponder(200, `[{"id": "test-repo", "url": "http://mock.repo/test-plugins.json"}]`))

	s := getArrayByteFromPath("test/metadata/valid-metadata.json", t)
	// mock metadata
	httpmock.RegisterResponder("GET", "http://mock.repo/test-plugins.json",
		httpmock.NewBytesResponder(200, s))

	repos := []string{
		"http://mock.repo/repository",
	}
	plugins := []plugin{
		{"Test.pluginA", "1.1.0"},
	}
	expectedResult := []CompatibilityResult{
		{"Test.pluginA", "1.1.0", false, "No compatible Spinnaker versions found for Plugin Test.pluginA@1.1.0"},
	}
	if got := ArePluginsCompatible("1.23.0", plugins, repos); !reflect.DeepEqual(got, expectedResult) {
		t.Errorf("method() = %v, want %v", got, expectedResult)
	}
}

func TestIncompatiblePluginWithoutConstrainForOlderSpinnakerVersion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock repository
	httpmock.RegisterResponder("GET", "http://mock.repo/repository",
		httpmock.NewStringResponder(200, `[{"id": "test-repo", "url": "http://mock.repo/test-plugins.json"}]`))

	s := getArrayByteFromPath("test/metadata/valid-metadata.json", t)
	// mock metadata
	httpmock.RegisterResponder("GET", "http://mock.repo/test-plugins.json",
		httpmock.NewBytesResponder(200, s))

	repos := []string{
		"http://mock.repo/repository",
	}
	plugins := []plugin{
		{"Test.pluginA", "1.1.0"},
	}
	expectedResult := []CompatibilityResult{
		{"Test.pluginA", "1.1.0", false, "No compatible Spinnaker versions found for Plugin Test.pluginA@1.1.0"},
	}
	if got := ArePluginsCompatible("1.19.7", plugins, repos); !reflect.DeepEqual(got, expectedResult) {
		t.Errorf("method() = %v, want %v", got, expectedResult)
	}
}

func TestMultiplePluginsWithCompatibleVersion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock repository
	httpmock.RegisterResponder("GET", "http://mock.repo/repository",
		httpmock.NewStringResponder(200, `[{"id": "test-repo", "url": "http://mock.repo/test-plugins.json"}]`))

	s := getArrayByteFromPath("test/metadata/valid-metadata.json", t)
	// mock metadata
	httpmock.RegisterResponder("GET", "http://mock.repo/test-plugins.json",
		httpmock.NewBytesResponder(200, s))

	repos := []string{
		"http://mock.repo/repository",
	}
	plugins := []plugin{
		{"Test.pluginA", "1.1.0"},
		{"Test.pluginB", "1.0.1"},
	}
	expectedResult := []CompatibilityResult{
		{"Test.pluginA", "1.1.0", true, "Plugin Test.pluginA@1.1.0 compatible with version 1.21.1"},
		{"Test.pluginB", "1.0.1", true, "Plugin Test.pluginB@1.0.1 compatible with version 1.21.0"},
	}
	if got := ArePluginsCompatible("1.21.0", plugins, repos); !reflect.DeepEqual(got, expectedResult) {
		t.Errorf("method() = %v, want %v", got, expectedResult)
	}
}

func TestMultiplePluginsWithIncompatibleVersion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock repository
	httpmock.RegisterResponder("GET", "http://mock.repo/repository",
		httpmock.NewStringResponder(200, `[{"id": "test-repo", "url": "http://mock.repo/test-plugins.json"}]`))

	s := getArrayByteFromPath("test/metadata/valid-metadata.json", t)
	// mock metadata
	httpmock.RegisterResponder("GET", "http://mock.repo/test-plugins.json",
		httpmock.NewBytesResponder(200, s))

	repos := []string{
		"http://mock.repo/repository",
	}
	plugins := []plugin{
		{"Test.pluginA", "1.1.0"},
		{"Test.pluginB", "1.0.1"},
	}
	expectedResult := []CompatibilityResult{
		{"Test.pluginA", "1.1.0", true, "Plugin Test.pluginA@1.1.0 compatible with version 1.20.0"},
		{"Test.pluginB", "1.0.1", false, "No compatible Spinnaker versions found for Plugin Test.pluginB@1.0.1"},
	}
	if got := ArePluginsCompatible("1.20.9", plugins, repos); !reflect.DeepEqual(got, expectedResult) {
		t.Errorf("method() = %v, want %v", got, expectedResult)
	}
}

func TestCompatiblePluginWhenCompatibilityMetadataIsMissing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock repository
	httpmock.RegisterResponder("GET", "http://mock.repo/repository",
		httpmock.NewStringResponder(200, `[{"id": "test-repo", "url": "http://mock.repo/test-plugins.json"}]`))

	s := getArrayByteFromPath("test/metadata/valid-metadata.json", t)
	// mock metadata
	httpmock.RegisterResponder("GET", "http://mock.repo/test-plugins.json",
		httpmock.NewBytesResponder(200, s))

	repos := []string{
		"http://mock.repo/repository",
	}
	plugins := []plugin{
		{"Test.pluginB", "1.0.0"},
	}
	expectedResult := []CompatibilityResult{
		{"Test.pluginB", "1.0.0", true, "Plugin Test.pluginB@1.0.0 does not contain compatibility constrain"},
	}
	if got := ArePluginsCompatible("1.20.0", plugins, repos); !reflect.DeepEqual(got, expectedResult) {
		t.Errorf("method() = %v, want %v", got, expectedResult)
	}
}

func TestCompatiblePluginWhenPluginDoesNotExist(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock repository
	httpmock.RegisterResponder("GET", "http://mock.repo/repository",
		httpmock.NewStringResponder(200, `[{"id": "test-repo", "url": "http://mock.repo/test-plugins.json"}]`))

	s := getArrayByteFromPath("test/metadata/valid-metadata.json", t)
	// mock metadata
	httpmock.RegisterResponder("GET", "http://mock.repo/test-plugins.json",
		httpmock.NewBytesResponder(200, s))

	repos := []string{
		"http://mock.repo/repository",
	}
	plugins := []plugin{
		{"Test.pluginC", "1.0.0"},
	}
	expectedResult := []CompatibilityResult{
		{"Test.pluginC", "1.0.0", true, "No releases found for Test.pluginC"},
	}
	if got := ArePluginsCompatible("1.20.0", plugins, repos); !reflect.DeepEqual(got, expectedResult) {
		t.Errorf("method() = %v, want %v", got, expectedResult)
	}
}
