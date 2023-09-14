package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/invocation_adapter/rest"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func statelesEgressAdapterFactory(adapterConfig model.StatelessEgressAdapterConfig, client ports.EgressClient) (ports.EgressAdapter, error) {
	switch adapterConfig.Variant {
	case constants.REST:
		return restEgressAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found")
		return nil, err
	}
}

func getOrCreateStatelessEgressClient(adapterConfig model.StatelessEgressAdapterConfig, clients *map[string]ports.EgressClient) ports.EgressClient {
	var client rest.RestClient
	switch adapterConfig.Variant {
	case constants.REST:
		key := string(adapterConfig.Variant)
		client, ok := (*clients)[key]
		if ok {
			logger.Logger.Infof("Client already exists reusing %v", key)
			return client
		}

		logger.Logger.Infof("created new client %s", key)
		client = rest.NewRestClient()
		(*clients)[key] = client
		return client
	default:
		logger.Logger.Fatalf("invalid protocol")
	}
	return client
}
