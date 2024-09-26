package validation

import (
	"testing"

	"github.com/hanapedia/hexagon/internal/config/infrastructure/yaml"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

// test multiple valid config
func TestServiceConfigsValidation(t *testing.T) {
	serviceUnitConfigs := []*model.ServiceUnitConfig{
		getServiceUnitConfig("./testdata/valid/service-a.yaml", t),
		getServiceUnitConfig("./testdata/valid/service-b.yaml", t),
		getServiceUnitConfig("./testdata/valid/service-c.yaml", t),
		getServiceUnitConfig("./testdata/valid/mongo.yaml", t),
	}
	errs := ValidateServiceUnitConfigs(serviceUnitConfigs)
	errs.Print()
	if errs.Exist() {
		t.Fail()
	}
}

func TestInvalidServiceMappingValidation(t *testing.T) {
	serviceUnitConfigs := []*model.ServiceUnitConfig{
		getServiceUnitConfig("./testdata/invalid/mapping/service-a.yaml", t),
		getServiceUnitConfig("./testdata/invalid/mapping/service-b.yaml", t),
		getServiceUnitConfig("./testdata/invalid/mapping/service-c.yaml", t),
		getServiceUnitConfig("./testdata/invalid/mapping/mongo.yaml", t),
	}
	errs := ValidateServiceUnitConfigs(serviceUnitConfigs)
	if !errs.Exist() {
		t.Fail()
	}
}

// test single valid config's fields
func TestServiceConfigValidation(t *testing.T) {
	serviceUnitConfig := getServiceUnitConfig("./testdata/valid/service-c.yaml", t)

	errs := ValidateServiceUnitConfigFields(serviceUnitConfig)
	errs.Print()
	if errs.Exist() {
		t.Fail()
	}
}

// test invalid service unit config fields
func TestInvalidServiceConfigFieldValidation(t *testing.T) {
	_, err := getInvalidServiceUnitConfig("./testdata/invalid/service_unit/invalid-adapter.yaml")

	if err == nil {
		t.Fail()
	}
}

func getServiceUnitConfig(path string, t *testing.T) *model.ServiceUnitConfig {
	serviceConfigLoader := yaml.YamlConfigLoader{Path: path}
	serviceConfig, err := serviceConfigLoader.Load()
	if err != nil {
		logger.Logger.Error(err)
		t.Fail()
	}
	return serviceConfig
}

func getInvalidServiceUnitConfig(path string) (*model.ServiceUnitConfig, error) {
	serviceConfigLoader := yaml.YamlConfigLoader{Path: path}
	serviceConfig, err := serviceConfigLoader.Load()
	return serviceConfig, err
}
