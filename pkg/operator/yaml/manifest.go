package yaml

import (
	k8syaml "sigs.k8s.io/yaml"
)

func GenerateManifest(object interface{}) []byte {
	yamlBytes, err := k8syaml.Marshal(object)
	if err != nil {
		panic(err)
	}
	return yamlBytes
}
