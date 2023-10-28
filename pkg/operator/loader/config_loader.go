package loader

import (
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

type ConfigLoader interface {
	Load() (model.ServiceUnitConfig, error)
}
