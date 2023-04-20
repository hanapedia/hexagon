package config

import (
	"reflect"
	"testing"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
)

func TestYamlConfigLoader_Load(t *testing.T) {
	testFilePath := "testdata/service-unit.yaml"

	flow1 := []core.Step{
		{AdapterId: "service-b.rest.read.getUserPreference"},
		{AdapterId: "service-c.rest.read.getUserHistory", Concurrent: true},
		{AdapterId: "service-d.rest.write.postUserLog"},
	}
	flow2 := []core.Step{
		{AdapterId: "service-c.rest.read.getUserHistory"},
		{AdapterId: "service-d.rest.write.updateUserPreference"},
	}
	handlerConfigs := []core.HandlerConfig{
		{Name: "getUser", Protocol: "rest", Action: "read", Steps: flow1},
		{Name: "updateUser", Protocol: "rest", Action: "write", Steps: flow2},
	}
	expectedConfig := core.ServiceUnitConfig{
		Name:           "service-a",
		HandlerConfigs: handlerConfigs,
	}

	ycl := YamlConfigLoader{
		Path: testFilePath,
	}

	config, err := ycl.Load()
	if err != nil {
		t.Fatalf("Error loading YAML file: %v", err)
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Loaded config does not match expected config.\nExpected: %+v\nActual: %+v", expectedConfig, config)
	}
}
