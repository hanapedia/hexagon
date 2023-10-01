package rest

import (
	"errors"
	"fmt"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hanapedia/the-bench/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/domain/contract"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/service-unit/utils"
)

// must implement ports.ServerAdapter
type RestServerAdapter struct {
	addr    string
	server  *fiber.App
	payload constants.PayloadSizeVariant
}

func NewRestServerAdapter(payload constants.PayloadSizeVariant) *RestServerAdapter {
	app := fiber.New()
	app.Use(fiber_logger.New())

	// enable tracing
	if config.GetEnvs().TRACING {
		app.Use(otelfiber.Middleware())
	}

	adapter := RestServerAdapter{addr: config.GetRestServerAddr(), server: app, payload: payload}

	return &adapter
}

func (rsa RestServerAdapter) Serve() error {
	return rsa.server.Listen(rsa.addr)
}

func (rsa RestServerAdapter) Register(serviceName string, handler *ports.PrimaryHandler) error {
	if handler.ServerConfig == nil {
		return errors.New(fmt.Sprintf("Invalid configuartion for handler %s.", handler.GetId()))
	}

	var err error
	switch handler.ServerConfig.Action {
	case "read":
		rsa.server.Get("/"+handler.ServerConfig.Route, func(c *fiber.Ctx) error {
			// call tasks
			errs := runtime.TaskSetHandler(c.Context(), handler.TaskSet)
			if errs != nil {
				for _, err := range errs {
					handler.LogTaskError(c.Context(), err)
				}

				restResponse := contract.RestResponseBody{
					Message: "Task failed",
				}
				return c.Status(fiber.StatusOK).JSON(restResponse)
			}

			// write response
			payload, err := utils.GeneratePayload(rsa.payload)
			if err != nil {
				return err
			}

			restResponse := contract.RestResponseBody{
				Message: fmt.Sprintf("Successfully ran %s, sending %s payload.", handler.GetId(), rsa.payload),
				Payload: &payload,
			}
			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	case "write":
		rsa.server.Post("/"+handler.ServerConfig.Route, func(c *fiber.Ctx) error {
			// call tasks
			errs := runtime.TaskSetHandler(c.Context(), handler.TaskSet)
			if errs != nil {
				for _, err := range errs {
					handler.LogTaskError(c.Context(), err)
				}
				restResponse := contract.RestResponseBody{
					Message: "Task failed",
				}
				return c.Status(fiber.StatusOK).JSON(restResponse)
			}

			// write response
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
