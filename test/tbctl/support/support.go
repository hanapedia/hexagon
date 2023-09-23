package support

import (
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/internal/tbctl/loader"
)

func GetServiceUnitConfig(path string) model.ServiceUnitConfig {
	return loader.GetConfig(path)
}
