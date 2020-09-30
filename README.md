
## Go Plugin Tools

This Module provides functionality for Spinnaker plugins with Go language

## How To Use This Package

In your Go files, simply use:

```
import "github.com/armory/plugin-tools/pkg/validation"
```

### Validations

The only exposed function is:

```
func ResolvePluginCompatibility(spinnaker SpinnakerVersion, plugins []plugin, repos []string) ([]CompatibilityResult, error)
```

This function return a list of `CompatibilityResult` with the verdict if a plugin is compatible with the given Spinnaker version or and error if something goes wrong when validating.

Here are the function arguments described:

* `SpinnakerVersion` This is the schema for a Spinnaker platform.
* `[]plugins` This is the list of plugins to check compatibility.
* `[]repos` This is the list of p4fj URL that contains the plugin's metadata.

The Spinnaker Version schema is composed by the SemVer of spinnaker, the kind of compatibility that we want to use and the service name.

The available kinds of compatibility are the following:
`spinnaker` for OSS.
`armory` for Armory Spinnaker Distribution.
`service` for Individual service compatibility.

The serviceName field is only required when using `service` kind.


All this information can be found in the Spinnaker configuration as shown bellow:
```yaml
      spinnaker:
        extensibility:
          plugins:
            Armory.CRDCheck:
              enabled: true
              version: 0.1.3
              extensions: {}
          repositories:
            crd-plugin-repository:
              url: https://raw.githubusercontent.com/armory-plugins/armory-crdcheck-plugin-releases/master/repositories.json
```

The Plugin metadata should look like the following:
```json
{
	"version": "1.1.14",
	"date": "2020-07-01T18:03:00.200Z",
	"requires": "orca>=0.0.0,deck>=0.0.0",
	"sha512sum":"f19deb40c2f386f1334a4ec6bf41bbb58296e489c37abcb80c93a5e423f2fb3522b45e8f9e5c7a188017c125b90bb0aea323e80f281fa1619a0ce769617e020e",
	"state": "",
	"compatibility": {
		"spinnaker" : ["1.21.1", "1.22.0"]
	},
        "url": "https://github.com/spinnaker-plugin-examples/pf4jStagePlugin/releases/download/v1.1.14/pf4jStagePlugin-v1.1.14.zip"
}
```

The validation will check for an exact match on the compatibility information, if an exact match was not find, it will check if there is a compatibility version that share the same minor version with the given spinnaker version. If any of these conditions are valid, it will mark the plugin as compatible.

If an error happened for a plugin (ex. plugin does not have compatibility information, plugin not found, etc) it will mark the plugin as compatible.