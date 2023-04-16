package core

import (
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

type ServiceUnit struct {
	Name           string
	ServerAdapters map[constants.ServerAdapterProtocol]*ServerAdapter
}
