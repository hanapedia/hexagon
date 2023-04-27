package tests

import (
	"log"
	"testing"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/config/yaml"
)

func TestConfigValidation(t *testing.T) {
	serviceUnitConfigs := []model.ServiceUnitConfig{
		getServiceUnitConfig("./testdata/service-a/service-unit.yaml"),
		getServiceUnitConfig("./testdata/service-b/service-unit.yaml"),
		getServiceUnitConfig("./testdata/service-c/service-unit.yaml"),
	}
	errs := model.ValidateConfig(serviceUnitConfigs)
	if len(errs.FieldErrors) > 0 {
		for _, fe := range errs.FieldErrors {
			log.Println(fe.Error())
		}
		t.Fail()
	}
	if len(errs.MappingErrors) > 0 {
		for _, me := range errs.MappingErrors {
			log.Println(me.Error())
		}
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
