package tests

import (
	"log"
	"testing"

	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/validation"
	"github.com/hanapedia/the-bench/pkg/operator/yaml"
)

// test multiple valid config
func TestServiceConfigsValidation(t *testing.T) {
	serviceUnitConfigs := []model.ServiceUnitConfig{
		getServiceUnitConfig("./testdata/valid/service-a.yaml"),
		getServiceUnitConfig("./testdata/valid/service-b.yaml"),
		getServiceUnitConfig("./testdata/valid/service-c.yaml"),
		getServiceUnitConfig("./testdata/valid/mongo.yaml"),
	}
	errs := validation.ValidateServiceUnitConfigs(serviceUnitConfigs)
	errs.Print()
	if errs.Exist() {
		t.Fail()
	}
}

func TestInvalidServiceMappingValidation(t *testing.T) {
	serviceUnitConfigs := []model.ServiceUnitConfig{
		getServiceUnitConfig("./testdata/invalid/mapping/service-a.yaml"),
		getServiceUnitConfig("./testdata/invalid/mapping/service-b.yaml"),
		getServiceUnitConfig("./testdata/invalid/mapping/service-c.yaml"),
		getServiceUnitConfig("./testdata/invalid/mapping/mongo.yaml"),
	}
	errs := validation.ValidateServiceUnitConfigs(serviceUnitConfigs)
	if !errs.Exist() {
		t.Fail()
	}
}

// test single valid config's fields
func TestServiceConfigValidation(t *testing.T) {
	serviceUnitConfig := getServiceUnitConfig("./testdata/valid/service-c.yaml")

	errs := validation.ValidateServiceUnitConfigFields(&serviceUnitConfig)
	errs.Print()
	if errs.Exist() {
		t.Fail()
	}
}

// test invalid service unit config fields
func TestInvalidServiceConfigFieldValidation(t *testing.T) {
	serviceUnitConfig := getServiceUnitConfig("./testdata/invalid/service_unit/invalidIngressAdapter.yaml")

	errs := validation.ValidateServiceUnitConfigFields(&serviceUnitConfig)
	if !errs.Exist() {
		t.Fail()
	}
}

// test invalid service unit config fields
func TestInvalidAdapterConfigValidation(t *testing.T) {
	serviceUnitConfig := getServiceUnitConfig("./testdata/invalid/adapter/egress-adapter.yaml")

	errs := validation.ValidateServiceUnitConfigFields(&serviceUnitConfig)
	if !errs.Exist() {
		t.Fail()
	}
}

// func TestMultipleStatefulServiceDefinitionValidation(t *testing.T) {
// 	serviceUnitConfig := getServiceUnitConfig("./testdata/invalid/service_unit/multipleStatefulAdapter.yaml")
//
// 	errs := validation.ValidateServiceUnitConfigFields(&serviceUnitConfig)
// 	errs.Print()
// 	if !errs.Exist() {
// 		t.Fail()
// 	}
// }
//
func getServiceUnitConfig(path string) model.ServiceUnitConfig {
	serviceConfigLoader := yaml.YamlConfigLoader{Path: path}
	serviceConfig, err := serviceConfigLoader.Load()
	if err != nil {
		log.Fatalf("Failed to load config from path %s. error: %s", path, err)
	}
	return serviceConfig
}
