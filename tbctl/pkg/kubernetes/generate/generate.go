package generate

import (
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

type ManifestGenerator struct {
	Input             string
	Output            string
	ServiceUnitConfig model.ServiceUnitConfig
}

func NewManifestGenerator(input string, output string) ManifestGenerator {
	return ManifestGenerator{
		Input:             input,
		Output:            output,
		ServiceUnitConfig: loader.GetConfig(input),
	}
}

func (mg ManifestGenerator) GenerateFromFile() {
	errs := mg.GenerateManifest()
	errs.Print()
	if errs.Exist() {
		logger.Logger.Fatal("Failed generating manifests.")
	}
}

func (mg ManifestGenerator) GenerateFromDirectory() {
	paths, err := loader.GetYAMLFiles(mg.Input)
	if err != nil {
		logger.Logger.Fatalf("Error reading from directory %s. %s", mg.Input, err)
	}

	for _, mg.Input = range paths {
		mg.GenerateFromFile()
	}

}

func (mg ManifestGenerator) GenerateManifest() ManifestErrors {
	var manfiestErrors ManifestErrors
	if hasStatefulAdapter(mg.ServiceUnitConfig.IngressAdapterConfigs) {
		manfiestErrors.Extend(mg.GenerateStatefulManifests())
		return manfiestErrors
	}
	brokerAdapters := getBrokerAdapters(mg.ServiceUnitConfig.IngressAdapterConfigs)
	if len(brokerAdapters) > 0 {
		for _, config := range brokerAdapters {
			manfiestErrors.Extend(mg.GenerateBrokerManifests(config))
		}
	}
	manfiestErrors.Extend(mg.GenerateStatelessManifests())
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

func getBrokerAdapters(ingressAdapterConfigs []model.IngressAdapterSpec) []model.BrokerIngressAdapterConfig {
	var configs []model.BrokerIngressAdapterConfig
	for _, ingresingressAdapterConfig := range ingressAdapterConfigs {
		if ingresingressAdapterConfig.BrokerIngressAdapterConfig != nil {
			configs = append(configs, *ingresingressAdapterConfig.BrokerIngressAdapterConfig)
		}
	}
	return configs
}
