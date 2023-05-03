package rest

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/logger"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/contract"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/common"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
)

// must implement core.ServerAdapter
type RestServerAdapter struct {
	addr   string
	server *fiber.App
}

func NewRestServerAdapter() RestServerAdapter {
	app := fiber.New()
	app.Use(fiber_logger.New())

	return RestServerAdapter{addr: config.GetRestServerAddr(), server: app}
}

func (rsa RestServerAdapter) Serve() error {
	return rsa.server.Listen(rsa.addr)
}

func (rsa RestServerAdapter) Register(handler *core.IngressAdapterHandler) error {
	if handler.StatelessIngressAdapterConfig == nil {
		return errors.New(fmt.Sprintf("Invalid configuartion for handler %s.", handler.GetId()))
	}

	var err error
	switch handler.StatelessIngressAdapterConfig.Action {
	case "read":
		rsa.server.Get("/"+handler.StatelessIngressAdapterConfig.Service, func(c *fiber.Ctx) error {
			egressAdapterErrors := common.TaskSetHandler(handler.TaskSets)
			for _, egressAdapterError := range egressAdapterErrors {
				logger.Logger.Errorf("Invocating %s failed: %s",
					reflect.TypeOf(egressAdapterError.EgressAdapter).Elem().Name(),
					egressAdapterError.Error,
				)
			}

			payload, err := utils.GenerateRandomString(constants.PayloadSize)
			if err != nil {
				return err
			}
			restResponse := contract.RestResponseBody{
				Message: fmt.Sprintf("Successfully ran %s, sending %vKB.", handler.GetId(), constants.PayloadSize),
				Payload: &payload,
			}
			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	case "write":
		rsa.server.Post("/"+handler.StatelessIngressAdapterConfig.Service, func(c *fiber.Ctx) error {
			egressAdapterErrors := common.TaskSetHandler(handler.TaskSets)
			for _, egressAdapterError := range egressAdapterErrors {
				logger.Logger.Errorf("Invocating %s failed: %s",
					reflect.TypeOf(egressAdapterError.EgressAdapter).Elem().Name(),
					egressAdapterError.Error,
				)
			}

			restResponse := contract.RestResponseBody{
				Message: fmt.Sprintf("Successfully ran %s.", handler.GetId()),
			}
			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	default:
		err = errors.New("Handler has no matching action")
	}
	return err
}
