package rest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/hanapedia/hexagon/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/internal/service-unit/domain/contract"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	util "github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

// must implement ports.PrimaryPort
type RestServerAdapter struct {
	addr        string
	server      *fiber.App
	payloadSize int64
}

func NewRestServerAdapter(payloadSize int64) *RestServerAdapter {
	app := fiber.New()

	// enable tracing
	if config.GetEnvs().TRACING {
		app.Use(otelfiber.Middleware())
	}

	adapter := RestServerAdapter{
		addr:        config.GetRestServerAddr(),
		server:      app,
		payloadSize: payloadSize,
	}

	return &adapter
}

func (rsa *RestServerAdapter) Serve() error {
	logger.Logger.Infof("Serving rest server at %s", rsa.addr)
	return rsa.server.Listen(rsa.addr)
}

func (rsa *RestServerAdapter) Register(handler *ports.PrimaryHandler) error {
	if handler.ServerConfig == nil {
		return errors.New(fmt.Sprintf("Invalid configuartion for handler %s.", handler.GetId()))
	}

	var err error
	switch handler.ServerConfig.Action {
	case constants.GET:
		rsa.server.Get("/"+handler.ServerConfig.Route, func(c *fiber.Ctx) error {
			// record time for logging
			startTime := time.Now()
			// defer log
			defer rsa.log(c.UserContext(), handler, startTime)

			// call tasks
			errs := runtime.TaskSetHandler(c.UserContext(), handler.TaskSet)
			if errs != nil {
				for _, err := range errs {
					handler.LogTaskError(c.UserContext(), err)
				}

				restResponse := contract.RestResponseBody{}
				return c.Status(fiber.StatusOK).JSON(restResponse)
			}

			// write response
			payload, err := util.GenerateRandomString(rsa.payloadSize)
			if err != nil {
				return err
			}

			restResponse := contract.RestResponseBody{
				Payload: &payload,
			}

			logger.Logger.Debugf("Ran %s, responding with %v bytes.", handler.GetId(), rsa.payloadSize)

			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	case constants.POST:
		rsa.server.Post("/"+handler.ServerConfig.Route, func(c *fiber.Ctx) error {
			// record time for logging
			startTime := time.Now()
			// defer log
			defer rsa.log(c.UserContext(), handler, startTime)

			// call tasks
			errs := runtime.TaskSetHandler(c.UserContext(), handler.TaskSet)
			if errs != nil {
				for _, err := range errs {
					handler.LogTaskError(c.UserContext(), err)
				}
				restResponse := contract.RestResponseBody{}
				return c.Status(fiber.StatusOK).JSON(restResponse)
			}

			// write response
			payload, err := util.GenerateRandomString(rsa.payloadSize)
			if err != nil {
				return err
			}

			// write response
			restResponse := contract.RestResponseBody{
				Payload: &payload,
			}

			logger.Logger.Debugf("Ran %s, responding with %v bytes.", handler.GetId(), rsa.payloadSize)

			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	default:
		err = errors.New("Handler has no matching action")
	}
	return err
}

func (rsa *RestServerAdapter) log(ctx context.Context, handler *ports.PrimaryHandler, startTime time.Time) {
	elapsed := time.Since(startTime).Milliseconds()
	unit := "ms"
	if elapsed == 0 {
		elapsed = time.Since(startTime).Microseconds()
		unit = "μs"
	}
	logger.Logger.WithContext(ctx).Infof("Invocation handled | %-40s | %10v %s", handler.GetId(), elapsed, unit)
}
