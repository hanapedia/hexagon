package generate

import (
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/defaults"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

type ManifestGenerator struct {
	// Input is the path to directory or file containing the servcie unit config yaml
	Input string

	// Output is the path to the output directory for the kubernetes manifests
	// NOTE: not file name. file names are automatically assigned by service name
	Output            string
	ServiceUnitConfig model.ServiceUnitConfig
}

func NewManifestGenerator(input, output string) ManifestGenerator {
	config := loader.GetConfig(input)
	defaults.SetDefauls(&config)

	return ManifestGenerator{
		Input:             input,
		Output:            output,
		ServiceUnitConfig: config,
	}
}

func (mg ManifestGenerator) GenerateFromFile() {
	errs := mg.GenerateManifest()
	errs.Print()
	if errs.Exist() {
		logger.Logger.Fatal("Failed generating manifests.")
	}
}

func GenerateFromDirectory(input, output string) {
	paths, err := loader.GetYAMLFiles(input)
	if err != nil {
		logger.Logger.Fatalf("Error reading from directory %s. %s", input, err)
	}

	for _, inputFile := range paths {
		mg := NewManifestGenerator(inputFile, output)
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
	if mg.ServiceUnitConfig.Gateway != nil {
		manfiestErrors.Extend(mg.GenerateLoadGeneratorManifests())
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
