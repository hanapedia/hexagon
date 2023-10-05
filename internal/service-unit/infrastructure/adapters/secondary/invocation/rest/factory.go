package rest

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
)

func RestInvocationAdapterFactory(adapterConfig *model.InvocationConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	var restAdapter ports.SecodaryPort
	var err error

	if restClient, ok := (client).(*RestClient); ok {
		port := config.GetEnvs().HTTP_PORT

		switch adapterConfig.Action {
		case constants.READ:
			restAdapter = &restReadAdapter{
				url:    fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route),
				client: restClient.Client,
			}
		case constants.WRITE:
			restAdapter = &restWriteAdapter{
				url:     fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route),
				client:  restClient.Client,
				payload: adapterConfig.Payload,
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
