package factory

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/invocation_adapter/rest"
)

func (egressAdapterDetails EgressAdapterDetails) restEgressAdapterFactory() (core.EgressAdapter, error) {
	var err error
	var restEgressAdapter core.EgressAdapter
	switch egressAdapterDetails.action {
	case "read":
		restEgressAdapter = rest.RestReadAdapter{URL: fmt.Sprintf("http://%s:8080/%s", egressAdapterDetails.serviceName, egressAdapterDetails.handlerName)}
	case "write":
		restEgressAdapter = rest.RestWriteAdapter{URL: fmt.Sprintf("http://%s:8080/%s", egressAdapterDetails.serviceName, egressAdapterDetails.handlerName)}
	default:
		err = errors.New("No matching protocol found")
	}
	return restEgressAdapter, err
}
