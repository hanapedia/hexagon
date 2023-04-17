package factory

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/invocation_adapter/rest"
)

func (invocationAdapterDetails InvocationAdapterDetails) restInvocationAdapterFactory() (core.InvocationAdapter, error) {
	var err error
	var invocationAdapter core.InvocationAdapter
	switch invocationAdapterDetails.action {
	case "read":
    invocationAdapter = rest.RestReadAdapter{URL: fmt.Sprintf("http://%s:8080/%s", invocationAdapterDetails.serviceName, invocationAdapterDetails.handlerName)}
	case "write":
		invocationAdapter = rest.RestWriteAdapter{URL: fmt.Sprintf("http://%s:8080/%s", invocationAdapterDetails.serviceName, invocationAdapterDetails.handlerName)}
	default:
		err = errors.New("No matching protocol found")
	}
	return invocationAdapter, err
}
