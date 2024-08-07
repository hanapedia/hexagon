package ports

import (
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

type ConfigLoader interface {
	Load() (*model.ServiceUnitConfig, error)
	LoadClusterConfig() (*model.ClusterConfig, error)
}
