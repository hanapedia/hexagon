package loader

import "github.com/hanapedia/the-bench/config/model"

type ConfigLoader interface {
	Load() (model.ServiceUnitConfig, error)
}
