package factory

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/invocation_adapter/rest"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

func (egressAdapterDetails EgressAdapterDetails) restEgressAdapterFactory() (core.EgressAdapter, error) {
	var err error
	var restEgressAdapter core.EgressAdapter
	switch egressAdapterDetails.action {
	case string(constants.READ):
		restEgressAdapter = rest.RestReadAdapter{URL: fmt.Sprintf("http://%s:8080/%s", egressAdapterDetails.serviceName, egressAdapterDetails.adapterName)}
	case string(constants.WRITE):
		restEgressAdapter = rest.RestWriteAdapter{URL: fmt.Sprintf("http://%s:8080/%s", egressAdapterDetails.serviceName, egressAdapterDetails.adapterName)}
	default:
		err = errors.New("No matching protocol found")
	}
	return restEgressAdapter, err
}
