package faultinjection
//
// import (
// 	"github.com/hanapedia/hexagon/pkg/operator/object/usecases"
// 	"github.com/hanapedia/hexagon/pkg/operator/yaml"
// )
//
// // ManifestGenerator generates manifest file for kafka topic
// func (mg ManifestGenerator) GenerateFaultInjectionManifests() ManifestErrors {
// 	// Open the manifestFile in append mode and with write-only permissions
// 	outPath := mg.getFilePath(mg.ServiceUnitConfig.Name, "fault-injection")
// 	manifestFile, err := createFile(outPath)
// 	if err != nil {
// 		return ManifestErrors{
// 			common: []CommonManifestError{
// 				NewCommonManifestError(mg.ServiceUnitConfig, "Unable to open output file."),
// 			},
// 		}
// 	}
// 	defer manifestFile.Close()
//
// 	manifest := usecases.CreateNetworkDelay(mg.ServiceUnitConfig.Name)
// 	manifestYAML := yaml.GenerateManifest(manifest)
// 	_, err = manifestFile.WriteString(formatManifest(manifestYAML))
// 	if err != nil {
// 		return ManifestErrors{
// 			common: []CommonManifestError{
// 				NewCommonManifestError(mg.ServiceUnitConfig, "Failed to write fault injection manifest"),
// 			},
// 		}
// 	}
// 	return ManifestErrors{}
// }
