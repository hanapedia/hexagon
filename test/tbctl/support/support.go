package support

import (
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/internal/tbctl/loader"
)

func GetServiceUnitConfig(path string) model.ServiceUnitConfig {
	return loader.GetConfig(path)
}
