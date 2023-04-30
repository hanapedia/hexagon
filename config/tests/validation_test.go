package tests

import (
	"log"
	"testing"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/config/yaml"
)

func TestServiceConfigsValidation(t *testing.T) {
	serviceUnitConfigs := []model.ServiceUnitConfig{
		getServiceUnitConfig("./testdata/service-a.yaml"),
		getServiceUnitConfig("./testdata/service-b.yaml"),
		getServiceUnitConfig("./testdata/service-c.yaml"),
	}
	errs := model.ValidateServiceUnitConfigs(serviceUnitConfigs)
	errs.Print()
	if len(errs.ServiceUnitFieldErrors) > 0 {
		t.Fail()
	}
	if len(errs.AdapterFieldErrors) > 0 {
		t.Fail()
	}
	if len(errs.MappingErrors) > 0 {
		t.Fail()
	}
}

func TestServiceConfigValidation(t *testing.T) {
	serviceUnitConfig := getServiceUnitConfig("./testdata/service-c.yaml")

	errs := model.ValidateServiceUnitConfigFields(serviceUnitConfig)
	errs.Print()
	if len(errs.ServiceUnitFieldErrors) > 0 {
		t.Fail()
	}
	if len(errs.AdapterFieldErrors) > 0 {
		t.Fail()
	}
	if len(errs.MappingErrors) > 0 {
		t.Fail()
	}
}

func TestInvalidServiceConfigValidation(t *testing.T) {
	serviceUnitConfig := getServiceUnitConfig("./testdata/invalid/invalidIngressAdapter.yaml")

	errs := model.ValidateServiceUnitConfigFields(serviceUnitConfig)
	errs.Print()
	if len(errs.ServiceUnitFieldErrors) > 0 {
		t.Fail()
	}
	if len(errs.AdapterFieldErrors) > 0 {
		t.Fail()
	}
	if len(errs.MappingErrors) > 0 {
		t.Fail()
	}
}

func getServiceUnitConfig(path string) model.ServiceUnitConfig {
	serviceConfigLoader := yaml.YamlConfigLoader{Path: path}
	serviceConfig, err := serviceConfigLoader.Load()
	if err != nil {
		log.Fatalf("Failed to load config from path %s. error: %s", path, err)
	}
	return serviceConfig
}
