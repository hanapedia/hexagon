package generate

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
)

func GenerateFromFile(input string, output string) {
	serviceUnitConfig := loader.GetConfig(input)
	errs := GenerateManifest(input, output, serviceUnitConfig)
    errs.Print()
	if errs.Exist() {
		logger.Logger.Fatal("Failed generating manifests.")
	}
}

func GenerateFromDirectory(input string, output string) {
	paths, err := loader.GetYAMLFiles(input)
	if err != nil {
		logger.Logger.Fatalf("Error reading from directory %s. %s", input, err)
	}

	for _, input = range paths {
		GenerateFromFile(input, output)
	}

}

func GenerateManifest(input string, output string, serviceUnitConfig model.ServiceUnitConfig) ManifestErrors {
	var manfiestErrors ManifestErrors
	if hasStatefulAdapter(serviceUnitConfig.IngressAdapterConfigs) {
		manfiestErrors.Extend(GenerateStatefulManifests(output, serviceUnitConfig))
		return manfiestErrors
	}
	if hasBrokerAdapter(serviceUnitConfig.IngressAdapterConfigs) {
		manfiestErrors.Extend(GenerateBrokerManifests(output, serviceUnitConfig))
	}
	manfiestErrors.Extend(GenerateStatelessManifests(input, output, serviceUnitConfig))
	return manfiestErrors
}

func hasStatefulAdapter(ingressAdapterConfigs []model.IngressAdapterSpec) bool {
	for _, ingresingressAdapterConfig := range ingressAdapterConfigs {
		if ingresingressAdapterConfig.StatefulIngressAdapterConfig != nil {
			return true
		}
	}
	return false
}

func hasBrokerAdapter(ingressAdapterConfigs []model.IngressAdapterSpec) bool {
	for _, ingresingressAdapterConfig := range ingressAdapterConfigs {
		if ingresingressAdapterConfig.BrokerIngressAdapterConfig != nil {
			return true
		}
	}
	return false
}
