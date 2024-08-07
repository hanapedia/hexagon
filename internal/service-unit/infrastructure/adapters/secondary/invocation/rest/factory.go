package rest

import (
	"errors"
	"fmt"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func RestInvocationAdapterFactory(adapterConfig *model.InvocationConfig, client secondary.SecondaryAdapterClient) (secondary.SecodaryPort, error) {
	var restAdapter secondary.SecodaryPort
	var err error

	if restClient, ok := (client).(*RestClient); ok {
		port := config.GetEnvs().HTTP_PORT

		switch adapterConfig.Action {
		case constants.GET, constants.READ:
			restAdapter = &restReadAdapter{
				url:    fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route),
				client: restClient.Client,
			}
		case constants.POST, constants.WRITE:
			restAdapter = &restWriteAdapter{
				url:         fmt.Sprintf("http://%s:%s/%s", adapterConfig.Service, port, adapterConfig.Route),
				client:      restClient.Client,
				payloadSize: model.GetPayloadSize(adapterConfig.Payload),
			}
		default:
			err = errors.New("No matching protocol found when creating rest secondary adapter.")
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	restAdapter.SetDestId(adapterConfig.GetId())

	logger.Logger.Debugf("Initialized Rest invocation adapter: %s", adapterConfig.GetId())
	return restAdapter, err
}
