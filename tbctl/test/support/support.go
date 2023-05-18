package support

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
)

func GetServiceUnitConfig(path string) model.ServiceUnitConfig {
	return loader.GetConfig(path)
}
