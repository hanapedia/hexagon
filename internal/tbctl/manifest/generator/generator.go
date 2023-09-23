package generator

import (
	"github.com/hanapedia/the-bench/internal/tbctl/loader"
	"github.com/hanapedia/the-bench/internal/tbctl/manifest/core"
	"github.com/hanapedia/the-bench/internal/tbctl/manifest/generator/broker"
	"github.com/hanapedia/the-bench/internal/tbctl/manifest/generator/loadgenerator"
	"github.com/hanapedia/the-bench/internal/tbctl/manifest/generator/serviceunit"
	"github.com/hanapedia/the-bench/internal/tbctl/manifest/generator/statefulunit"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
)

type ManifestGenerator struct {
	// Input is the path to directory or file containing the servcie unit config yaml
	Input string

	// Output is the path to the output directory for the kubernetes manifests
	// NOTE: not file name. file names are automatically assigned by service name
	Output            string
	ServiceUnitConfig *model.ServiceUnitConfig
}

func NewManifestGenerator(input, output string) ManifestGenerator {
	config := loader.GetConfig(input)

	return ManifestGenerator{
		Input:             input,
		Output:            output,
		ServiceUnitConfig: &config,
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

func (mg ManifestGenerator) GenerateManifest() core.ManifestErrors {
	var manfiestErrors core.ManifestErrors

	// generate stateful unit manifest and exit if the config is for stateful unit
	if core.HasRepositoryAdapter(mg.ServiceUnitConfig) {
		statefuleManifest := statefulunit.NewStatefulUnitManifest(mg.ServiceUnitConfig)
		path := core.GetFilePath(mg.Output, mg.ServiceUnitConfig.Name, "stateful-unit")
		manfiestErrors.Extend(statefuleManifest.Generate(mg.ServiceUnitConfig, path))
		return manfiestErrors
	}

	// generate broker manifest if the config has consumer adapter
	if core.HasConsumerAdapters(mg.ServiceUnitConfig) {
		brokerManifest := broker.NewBrokerManifest(mg.ServiceUnitConfig)
		path := core.GetFilePath(mg.Output, mg.ServiceUnitConfig.Name, "broker")
		manfiestErrors.Extend(brokerManifest.Generate(mg.ServiceUnitConfig, path))
	}

	// generate loadgenerator manifest if the config has loadgenerator
	if core.HasGatewayConfig(mg.ServiceUnitConfig) {
		loadGeneratorManifest := loadgenerator.NewLoadGeneratorManifest(mg.ServiceUnitConfig)
		path := core.GetFilePath(mg.Output, mg.ServiceUnitConfig.Name, "load-generator")
		manfiestErrors.Extend(loadGeneratorManifest.Generate(mg.ServiceUnitConfig, path))
	}

	serviceUnitManifest := serviceunit.NewServiceUnitManifest(mg.ServiceUnitConfig, mg.Input)
	path := core.GetFilePath(mg.Output, mg.ServiceUnitConfig.Name, "service-unit")
	manfiestErrors.Extend(serviceUnitManifest.Generate(mg.ServiceUnitConfig, path))

	return manfiestErrors
}
