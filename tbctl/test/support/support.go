package support

import (
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
)

func GetServiceUnitConfig(path string) model.ServiceUnitConfig {
	return loader.GetConfig(path)
}
