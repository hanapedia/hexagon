package loader

import (
	model "github.com/hanapedia/the-bench/pkg/api/v1"
)

type ConfigLoader interface {
	Load() (model.ServiceUnitConfig, error)
}
