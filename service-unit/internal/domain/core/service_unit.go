package core

import (
	"github.com/hanapedia/the-bench/service-unit/pkg/shared"
)

type ServiceUnit struct {
	Name           string
	ServerAdapters map[shared.ServerAdapterProtocol]*ServerAdapter
}
