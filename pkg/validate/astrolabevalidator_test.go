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
	verdict, err := v.IsPluginCompatibleWithPlatform("1.23.1", "Armory.Test", "1.0.0", repos)
	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if verdict != Compatible {
		t.Errorf("method() = %s but want = %s", verdict, Compatible)
	}
}

func TestIncompatiblePluginWithPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	verdict, err := v.IsPluginCompatibleWithPlatform("1.22.5", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if verdict != NotCompatible {
		t.Errorf("method() = %s but want = %s", verdict, NotCompatible)
	}
}

func TestIncompatiblePluginThatOperateMultipleServicesWithPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	verdict, err := v.IsPluginCompatibleWithPlatform("1.21.7", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if verdict != NotCompatible {
		t.Errorf("method() = %s but want = %s", verdict, NotCompatible)
	}
}

func TestCompatiblePluginWithCustomPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	verdict, err := v.IsPluginCompatibleWithPlatform("2.22.1", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if verdict != Compatible {
		t.Errorf("method() = %s but want = %s", verdict, Compatible)
	}
}

func TestIncompatiblePluginWithCustomPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	verdict, err := v.IsPluginCompatibleWithPlatform("2.22.1", "Armory.Test", "0.1.7", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if verdict != NotCompatible {
		t.Errorf("method() = %s but want = %s", verdict, NotCompatible)
	}
}

func TestIncompatiblePluginThatOperatesMultipleServicesWithCustomPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	verdict, err := v.IsPluginCompatibleWithPlatform("2.22.0", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if verdict != NotCompatible {
		t.Errorf("method() = %s but want = %s", verdict, NotCompatible)
	}
}

func TestReturnErrorWithUnknownPlatform(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	if verdict, err := v.IsPluginCompatibleWithPlatform("fake-platform", "Armory.Test", "1.0.0", repos); err == nil || verdict != Unknown {
		t.Errorf("An error was expected and method() = %s but want = %s", verdict, Unknown)
	}
}

func TestReturnErrorWhenMetadataIsMissing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewStringResponder(404, ""))

	v := NewAstrolabeValidator()
	if verdict, err := v.IsPluginCompatibleWithPlatform("1.23.1", "Armory.Test", "1.0.0", repos); err == nil || verdict != Unknown {
		t.Errorf("An error was expected and method() = %s but want = %s", verdict, Unknown)
	}
}

func TestReturnErrorWhenReleaseIsMissing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	if verdict, err := v.IsPluginCompatibleWithPlatform("1.23.1", "Armory.Test", "2.0.0", repos); err == nil || verdict != Unknown {
		t.Errorf("An error was expected and method() = %s but want = %s", verdict, Unknown)
	}
}

func TestReturnErrorWhenPlatformIsMissing(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	if verdict, err := v.IsPluginCompatibleWithPlatform("1.20.0", "Armory.Test", "1.0.0", repos); err == nil || verdict != Unknown {
		t.Errorf("An error was expected and method() = %s but want = %s", verdict, Unknown)
	}
}

func TestCompatiblePluginWithService(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	verdict, err := v.IsPluginCompatibleWithService("orca", "5.0.0", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if verdict != Compatible {
		t.Errorf("method() = %s but want = %s", verdict, Compatible)
	}
}

func TestIncompatiblePluginWithService(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	verdict, err := v.IsPluginCompatibleWithService("orca", "4.0.0", "Armory.Test", "1.0.0", repos)

	if err != nil {
		t.Errorf("method() return error: %s", err)
	}
	if verdict != NotCompatible {
		t.Errorf("method() = %s but want = %s", verdict, Compatible)
	}
}

func TestReturnErrorWithUnknownService(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := getArrayByteFromPath("test/metadata/compatibility-metadata-test.json", t)
	httpmock.RegisterResponder("GET", "http://mock.repo/compatibility/Armory.Test.json",
		httpmock.NewBytesResponder(200, s))

	v := NewAstrolabeValidator()
	if verdict, err := v.IsPluginCompatibleWithService("fake-service", "fake-version", "Armory.Test", "1.0.0", repos); err == nil || verdict != Unknown {
		t.Errorf("An error was expected and method() = %s but want = %s", verdict, Unknown)
	}
}
