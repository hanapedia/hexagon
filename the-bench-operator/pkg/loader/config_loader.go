package loader

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

type ConfigLoader interface {
	Load() (model.ServiceUnitConfig, error)
}
