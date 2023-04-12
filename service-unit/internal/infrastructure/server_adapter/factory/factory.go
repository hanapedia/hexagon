package factory

import (
	"log"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/server_adapter/rest"
	"github.com/hanapedia/the-bench/service-unit/pkg/shared"
)

func NewServerAdapter(serverAdapterProtocol shared.ServerAdapterProtocol) *core.ServerAdapter {
	var serverAdapter core.ServerAdapter

	switch serverAdapterProtocol {
	case "rest":
		serverAdapter = rest.NewRestServerAdapter()
	default:
		log.Fatal("Adapter currently unsupported.")
	}

	return &serverAdapter
}

// Takes the pointer to the slice of ServerAdapters
// Update or insert ServiceAdapter based on the handler input.
// Does not return any value
func RegiserHandlerToServerAdapter(serverAdapterProtocol shared.ServerAdapterProtocol, serverAdapter *core.ServerAdapter, handler *core.Handler) error {
	err := (*serverAdapter).Register(handler)

	return err
}
