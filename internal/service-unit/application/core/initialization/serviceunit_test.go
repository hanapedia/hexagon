package initialization_test

import (
	"testing"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/initialization"
	ports "github.com/hanapedia/hexagon/internal/config/application/ports"
	"github.com/hanapedia/hexagon/internal/config/infrastructure/yaml"
)

var testdataDir string = "../../../../testdata/config"

func newConfigLoader(path string) ports.ConfigLoader {
	return yaml.YamlConfigLoader{Path: path}
}

func setupServiceUnit(testDir string) {
	configLoader := newConfigLoader(testdataDir + testDir)
	serviceUnitConfig := initialization.GetConfig(configLoader)

	// init telemetry
	initialization.InitTracing(serviceUnitConfig.Name)

	serviceUnit := initialization.NewServiceUnit(serviceUnitConfig)

	// setup service unit
	serviceUnit.Setup()
}

func TestValidServiceUnitSetup(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("The code panicked with %v", r)
        }
    }()

	testDirs := []string{
		"/tracing/gateway/service-unit.yaml",
		"/tracing/chain1/service-unit.yaml",
		"/kafka/producer/service-unit.yaml",
		"/kafka/consumer/service-unit.yaml",
		"/mongo/client/service-unit.yaml",
	}

	for _, testDir := range testDirs {
		setupServiceUnit(testDir)
	}
}

// func TestInvalidServiceUnitSetup(t *testing.T) {
// 	defer func() {
// 		r := recover()
// 		if r == nil {
// 			t.Fatal("The code did not panic")
// 		}
//
// 		if r != "Invalid primary adapter config" {
// 			t.Fatalf("Expected a different panic value. Got %v", r)
// 		}
// 	}()
// 	setupServiceUnit("/mongo/db/service-unit.yaml")
// }
