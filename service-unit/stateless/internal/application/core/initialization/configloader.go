package initialization

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/defaults"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/loader"
	l "github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/validation"
)

// GetConfig parses service-unit config via given formant
func GetConfig(configLoader loader.ConfigLoader) model.ServiceUnitConfig {
	config, err := configLoader.Load()
	if err != nil {
		l.Logger.Fatalf("Error loading config: %v", err)
	}

	defaults.SetDefaults(&config)

	errs := validation.ValidateServiceUnitConfigFields(&config)
	if errs.Exist() {
		errs.Print()
		l.Logger.Fatalln("Validation failed. Aborted.")
	}
	return config
}
