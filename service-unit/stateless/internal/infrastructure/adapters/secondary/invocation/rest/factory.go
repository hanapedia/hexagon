package rest

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
)

func RestInvocationAdapterFactory(adapterConfig model.InvocationConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	var restAdapter ports.SecodaryPort
	var err error

	if restClient, ok := (client).(RestClient); ok {
		port := config.GetEnvs().HTTP_PORT

		switch adapterConfig.Action {
		case constants.READ:
			restAdapter = &RestReadAdapter{
				URL:    fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route),
				Client: restClient.Client,
			}
		case constants.WRITE:
			restAdapter = &RestWriteAdapter{
				URL:    fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route),
				Client: restClient.Client,
			}
		default:
			err = errors.New("No matching protocol found when creating rest secondary adapter.")
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	restAdapter.SetDestId(adapterConfig.GetId())

	return restAdapter, err
}
