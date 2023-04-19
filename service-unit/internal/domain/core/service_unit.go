package core

import (
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

type ServiceUnit struct {
	Name           string
	Config         *ConfigLoader
	ServerAdapters *map[constants.AdapterProtocol]*IngressAdapter
}
