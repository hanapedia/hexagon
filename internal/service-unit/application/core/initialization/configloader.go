package initialization

import (
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/loader"
	l "github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/validation"
)

// GetConfig parses service-unit config via given formant
func GetConfig(configLoader loader.ConfigLoader) model.ServiceUnitConfig {
	config, err := configLoader.Load()
	if err != nil {
		l.Logger.Fatalf("Error loading config: %v", err)
	}

	errs := validation.ValidateServiceUnitConfigFields(&config)
	if errs.Exist() {
		errs.Print()
		l.Logger.Fatalln("Validation failed. Aborted.")
	}
	return config
}
