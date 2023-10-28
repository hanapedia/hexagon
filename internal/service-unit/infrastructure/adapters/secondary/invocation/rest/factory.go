package rest

import (
	"errors"
	"fmt"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
)

func RestInvocationAdapterFactory(adapterConfig *model.InvocationConfig, client ports.SecondaryAdapterClient) (ports.SecodaryPort, error) {
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
