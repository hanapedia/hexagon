package tests

import (
	"testing"

	"github.com/hanapedia/the-bench/internal/tbctl/generate"
)

func TestGenearateStatelessManifest(t *testing.T) {
	manifestGenerator := generate.NewManifestGenerator("./data/stateless.yaml", "./data/output/stateless.yaml")
	errs := manifestGenerator.GenerateStatelessManifests()
	if errs.Exist() {
		errs.Print()
		t.Fatal("Failed to generate manifests.")
	}
}

func TestGenearateMongoManifest(t *testing.T) {
	manifestGenerator := generate.NewManifestGenerator("./data/mongo.yaml", "./data/output/stateful.yaml")
	errs := manifestGenerator.GenerateStatefulManifests()
	if errs.Exist() {
		errs.Print()
		t.Fatal("Failed to generate manifests.")
	}
}

func TestGenearateKafkaTopic(t *testing.T) {
	manifestGenerator := generate.NewManifestGenerator("./data/kafka.yaml", "./data/output/kafka.yaml")
	if manifestGenerator.ServiceUnitConfig.AdapterConfigs[0].ConsumerConfig != nil {
		config := manifestGenerator.ServiceUnitConfig.AdapterConfigs[0].ConsumerConfig
		errs := manifestGenerator.GenerateBrokerManifests(*config)
		if errs.Exist() {
			errs.Print()
			t.Fatal("Failed to generate manifests.")
		}
	} else {
		t.Fail()
	}
}

func TestGenearateLoadGenerator(t *testing.T) {
	manifestGenerator := generate.NewManifestGenerator("./data/load_generator.yaml", "./data/output/load_generator.yaml")
	errs := manifestGenerator.GenerateLoadGeneratorManifests()
	if errs.Exist() {
		errs.Print()
		t.Fatal("Failed to generate manifests.")
	}
}
