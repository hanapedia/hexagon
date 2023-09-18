package generate

import (
	"github.com/hanapedia/the-bench/internal/tbctl/loader"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/defaults"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
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
	defaults.SetDefaults(&config)

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
	if hasRepositoryAdapter(mg.ServiceUnitConfig.AdapterConfigs) {
		manfiestErrors.Extend(mg.GenerateStatefulManifests())
		return manfiestErrors
	}
	brokerAdapters := getBrokerAdapters(mg.ServiceUnitConfig.AdapterConfigs)
	if len(brokerAdapters) > 0 {
		for _, config := range brokerAdapters {
			manfiestErrors.Extend(mg.GenerateBrokerManifests(config))
		}
	}
	if mg.ServiceUnitConfig.Gateway != nil {
		manfiestErrors.Extend(mg.GenerateLoadGeneratorManifests())
	}

	// generate stateless manifest
	manfiestErrors.Extend(mg.GenerateStatelessManifests())

	// generate namespace
	// manfiestErrors.Extend(mg.GenerateNamespaceManifests())

	// generate fault injection manifest
	// manfiestErrors.Extend(mg.GenerateFaultInjectionManifests())
	return manfiestErrors
}

func hasRepositoryAdapter(primaryAdapterConfigs []model.PrimaryAdapterSpec) bool {
	for _, primaryAdapterConfig := range primaryAdapterConfigs {
		if primaryAdapterConfig.RepositoryConfig != nil {
			return true
		}
	}
	return false
}

func getBrokerAdapters(primaryAdapterConfigs []model.PrimaryAdapterSpec) []model.ConsumerConfig {
	var configs []model.ConsumerConfig
	for _, primaryAdapterConfig := range primaryAdapterConfigs {
		if primaryAdapterConfig.ConsumerConfig != nil {
			configs = append(configs, *primaryAdapterConfig.ConsumerConfig)
		}
	}
	return configs
}
