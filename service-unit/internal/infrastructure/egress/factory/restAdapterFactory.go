package factory

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/invocation_adapter/rest"
)

func restEgressAdapterFactory(adapterConfig model.StatelessEgressConfig) (core.EgressAdapter, error) {
	var err error
	var restEgressAdapter core.EgressAdapter
	switch adapterConfig.Action {
	case constants.READ:
		restEgressAdapter = rest.RestReadAdapter{URL: fmt.Sprintf("http://%s:8080/%s", adapterConfig.Service, adapterConfig.Route)}
	case constants.WRITE:
		restEgressAdapter = rest.RestWriteAdapter{URL: fmt.Sprintf("http://%s:8080/%s", adapterConfig.Service, adapterConfig.Route)}
	default:
		err = errors.New("No matching protocol found")
	}
	return restEgressAdapter, err
}
