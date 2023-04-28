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
	if len(errs.ServiceUnitFieldErrors) > 0 {
		for _, fe := range errs.ServiceUnitFieldErrors {
			log.Println(fe.Error())
		}
		t.Fail()
	}
	if len(errs.AdapterFieldErrors) > 0 {
		for _, fe := range errs.AdapterFieldErrors {
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

func TestServiceConfigValidation(t *testing.T) {
	serviceUnitConfig := getServiceUnitConfig("./testdata/service-c.yaml")
	
	sufe, afe := model.ValidateServiceUnitConfigFields(serviceUnitConfig)
	if len(sufe) > 0 {
		for _, fe := range sufe {
			log.Println(fe.Error())
		}
		t.Fail()
	}
	if len(afe) > 0 {
		for _, fe := range afe {
			log.Println(fe.Error())
		}
		t.Fail()
	}
}

func TestInvalidServiceConfigValidation(t *testing.T) {
	serviceUnitConfig := getServiceUnitConfig("./testdata/invalid/invalidIngressAdapter.yaml")
	
	sufe, afe := model.ValidateServiceUnitConfigFields(serviceUnitConfig)
	if len(sufe) > 0 {
		for _, fe := range sufe {
			log.Println(fe.Error())
		}
		t.Fail()
	}
	if len(afe) > 0 {
		for _, fe := range afe {
			log.Println(fe.Error())
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
