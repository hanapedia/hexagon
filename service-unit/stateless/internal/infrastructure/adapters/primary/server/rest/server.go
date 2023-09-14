package rest

import (
	"errors"
	"fmt"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/contract"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/core/runtime"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
)

// must implement ports.ServerAdapter
type RestServerAdapter struct {
	addr   string
	server *fiber.App
}

func NewRestServerAdapter() RestServerAdapter {
	app := fiber.New()
	app.Use(fiber_logger.New())

	// enable tracing
	if config.GetEnvs().TRACING {
		app.Use(otelfiber.Middleware())
	}

	return RestServerAdapter{addr: config.GetRestServerAddr(), server: app}
}

func (rsa RestServerAdapter) Serve() error {
	return rsa.server.Listen(rsa.addr)
}

func (rsa RestServerAdapter) Register(serviceName string, handler *ports.PrimaryHandler) error {
	if handler.StatelessPrimaryAdapterConfig == nil {
		return errors.New(fmt.Sprintf("Invalid configuartion for handler %s.", handler.GetId(serviceName)))
	}

	var err error
	switch handler.StatelessPrimaryAdapterConfig.Action {
	case "read":
		rsa.server.Get("/"+handler.StatelessPrimaryAdapterConfig.Route, func(c *fiber.Ctx) error {
			// call tasks
			egressAdapterErrors := runtime.TaskSetHandler(c.Context(), handler.TaskSets)
			ports.LogSecondaryPortErrors(&egressAdapterErrors)


			// write response
			payload, err := utils.GenerateRandomString(constants.PayloadSize)
			if err != nil {
				return err
			}
			restResponse := contract.RestResponseBody{
				Message: fmt.Sprintf("Successfully ran %s, sending %vKB.", handler.GetId(serviceName), constants.PayloadSize),
				Payload: &payload,
			}
			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	case "write":
		rsa.server.Post("/"+handler.StatelessPrimaryAdapterConfig.Route, func(c *fiber.Ctx) error {
			// call tasks
			egressAdapterErrors := runtime.TaskSetHandler(c.Context() ,handler.TaskSets)
			ports.LogSecondaryPortErrors(&egressAdapterErrors)

			// write response
			restResponse := contract.RestResponseBody{
				Message: fmt.Sprintf("Successfully ran %s.", handler.GetId(serviceName)),
			}
			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	default:
		err = errors.New("Handler has no matching action")
	}
	return err
}
