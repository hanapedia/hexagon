package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
)

// Takes the pointer to the slice of ServerAdapters
// Update or insert ServiceAdapter based on the handler input.
// Does not return any value
func UpsertServerAdapter(serverAdapters *[]*core.ServerAdapter, handler *core.Handler) error {
	var err error = nil

	switch handler.Protocol {
	case "rest":
        RestServerAdapterFactory(serverAdapters, handler)
	default:
		err = errors.New("Handler has no matching protocol")
	}

	return err
}
