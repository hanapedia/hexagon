package namespace
//
// import (
// 	"github.com/hanapedia/hexagon/pkg/operator/object/usecases"
// 	"github.com/hanapedia/hexagon/pkg/operator/yaml"
// )
//
// // GenerateNamespaceManifests generates namespace manifest
// func (mg ManifestGenerator) GenerateNamespaceManifests() ManifestErrors {
// 	// Open the manifestFile in append mode and with write-only permissions
// 	outPath := mg.getFilePath(mg.ServiceUnitConfig.Name, "namespace")
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
// 	namespace := usecases.CreateNamespace()
// 	namespaceYAML := yaml.GenerateManifest(namespace)
// 	_, err = manifestFile.WriteString(formatManifest(namespaceYAML))
// 	if err != nil {
// 		return ManifestErrors{
// 			common: []CommonManifestError{
// 				NewCommonManifestError(mg.ServiceUnitConfig, "Failed to write namespace manifest"),
// 			},
// 		}
// 	}
// 	return ManifestErrors{}
// }
