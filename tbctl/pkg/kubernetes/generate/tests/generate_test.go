package tests

import (
	"testing"

	"github.com/hanapedia/the-bench/tbctl/pkg/kubernetes/generate"
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
	if manifestGenerator.ServiceUnitConfig.IngressAdapterConfigs[0].BrokerIngressAdapterConfig != nil {
		config := manifestGenerator.ServiceUnitConfig.IngressAdapterConfigs[0].BrokerIngressAdapterConfig
		errs := manifestGenerator.GenerateBrokerManifests(*config)
		if errs.Exist() {
			errs.Print()
			t.Fatal("Failed to generate manifests.")
		}
	} else {
		t.Fail()
	}
}
