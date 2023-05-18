package factory

import (
	"errors"
	"fmt"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/invocation_adapter/rest"
)

func restEgressAdapterFactory(adapterConfig model.StatelessEgressAdapterConfig) (core.EgressAdapter, error) {
	var err error
	var restEgressAdapter core.EgressAdapter

	port := config.GetEnvs().HTTP_PORT

	switch adapterConfig.Action {
	case constants.READ:
		restEgressAdapter = rest.RestReadAdapter{URL: fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route)}
	case constants.WRITE:
		restEgressAdapter = rest.RestWriteAdapter{URL: fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route)}
	default:
		err = errors.New("No matching protocol found when creating rest egress adapter.")
	}
	return restEgressAdapter, err
}
