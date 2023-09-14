package invocation

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/invocation/rest"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func NewSecondaryAdapter(adapterConfig model.StatelessEgressAdapterConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.REST:
		return rest.RestEgressAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found")
		return nil, err
	}
}

func GetOrCreateClient(adapterConfig model.StatelessEgressAdapterConfig, clients *map[string]ports.SecondaryAdapter) ports.SecondaryAdapter {
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
