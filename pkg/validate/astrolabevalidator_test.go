package validate

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"os"
	"testing"
)

var repos = []PluginRepository{
	{
		"mocked",
		"http://mock.repo/repositories.json",
	},
}

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

func TestCompatiblePluginWithPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()

	isValid, err := v.IsPluginCompatibleWithPlatform("1.23.1", "Armory.Test", "1.0.0", repos)
	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if !isValid {
		t.Errorf("method() = %t but want = %t", isValid, true)
	}
}

func TestIncompatiblePluginWithPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	isValid, err := v.IsPluginCompatibleWithPlatform("1.22.5", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if isValid {
		t.Errorf("method() = %t but want = %t", isValid, false)
	}
}

func TestCompatiblePluginWithCustomPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	isValid, err := v.IsPluginCompatibleWithPlatform("2.22.1", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if !isValid {
		t.Errorf("method() = %t but want = %t", isValid, true)
	}
}

func TestIncompatiblePluginWithCustomPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	isValid, err := v.IsPluginCompatibleWithPlatform("2.22.1", "Armory.Test", "0.1.7", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if isValid {
		t.Errorf("method() = %t but want = %t", isValid, false)
	}
}

func TestReturnErrorWithUnknownPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	if _, err := v.IsPluginCompatibleWithPlatform("fake-platform", "Armory.Test", "1.0.0", repos); err == nil {
		t.Errorf("An error was expected")
	}
}

func TestReturnErrorWhenMetadataIsMissing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewStringResponder(404, ""))

	v := NewAstrolabeValidator()
	if _, err := v.IsPluginCompatibleWithPlatform("1.23.1", "Armory.Test", "1.0.0", repos); err == nil {
		t.Errorf("An error was expected")
	}
}

func TestReturnErrorWhenReleaseIsMissing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	if _, err := v.IsPluginCompatibleWithPlatform("1.23.1", "Armory.Test", "2.0.0", repos); err == nil {
		t.Errorf("An error was expected")
	}
}

func TestReturnErrorWhenPlatformIsMissing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	if _, err := v.IsPluginCompatibleWithPlatform("1.20.0", "Armory.Test", "1.0.0", repos); err == nil {
		t.Errorf("An error was expected")
	}
}

func TestCompatiblePluginWithService(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	isValid, err := v.IsPluginCompatibleWithService("orca", "5.0.0", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if !isValid {
		t.Errorf("method() = %t but want = %t", isValid, true)
	}
}

func TestIncompatiblePluginWithService(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	isValid, err := v.IsPluginCompatibleWithService("orca", "4.0.0", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if isValid {
		t.Errorf("method() = %t but want = %t", isValid, false)
	}
}

func TestReturnErrorWithUnknownService(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	if _, err := v.IsPluginCompatibleWithService("fake-service", "fake-version", "Armory.Test", "1.0.0", repos); err == nil {
		t.Errorf("An error was expected")
	}
}
