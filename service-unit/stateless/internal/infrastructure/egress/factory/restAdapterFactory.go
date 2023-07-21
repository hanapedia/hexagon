package factory

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/invocation_adapter/rest"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
)

func restEgressAdapterFactory(adapterConfig model.StatelessEgressAdapterConfig, client core.EgressClient) (core.EgressAdapter, error) {
	var restEgressAdapter core.EgressAdapter
	var err error

	if restClient, ok := (client).(rest.RestClient); ok {
		port := config.GetEnvs().HTTP_PORT

		switch adapterConfig.Action {
		case constants.READ:
			restEgressAdapter = rest.RestReadAdapter{
				URL:    fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route),
				Client: restClient.Client,
			}
		case constants.WRITE:
			restEgressAdapter = rest.RestWriteAdapter{
				URL:    fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route),
				Client: restClient.Client,
			}
		default:
			err = errors.New("No matching protocol found when creating rest egress adapter.")
		}
	} else {
		err = errors.New("Unmatched client instance")
	}
	return restEgressAdapter, err
}
