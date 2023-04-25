package rest

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/contract"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/common"
	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
)

// must implement core.ServerAdapter
type RestServerAdapter struct {
	addr   string
	server *fiber.App
}

func NewRestServerAdapter() RestServerAdapter {
	app := fiber.New()
	app.Use(logger.New())

	return RestServerAdapter{addr: constants.RestServerAddr, server: app}
}

func (rsa RestServerAdapter) Serve() error {
	return rsa.server.Listen(rsa.addr)
}

func (rsa RestServerAdapter) Register(handler *core.Handler) error {
	var err error
	switch handler.Action {
	case "read":
		rsa.server.Get("/"+handler.Name, func(c *fiber.Ctx) error {
			egressAdapterErrors := common.TaskSetHandler(handler.TaskSets)
			for _, egressAdapterError := range egressAdapterErrors {
				log.Printf("Invocating %s failed: %s",
					reflect.TypeOf(egressAdapterError.EgressAdapter).Elem().Name(),
					egressAdapterError.Error,
				)
			}

			payload, err := utils.GenerateRandomString(constants.PayloadSize)
			if err != nil {
				return err
			}
			restResponse := contract.RestResponseBody{
				Message: fmt.Sprintf("Successfully ran %s, sending %vKB.", handler.ID, constants.PayloadSize),
				Payload: &payload,
			}
			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	case "write":
		rsa.server.Post("/"+handler.Name, func(c *fiber.Ctx) error {
			egressAdapterErrors := common.TaskSetHandler(handler.TaskSets)
			for _, egressAdapterError := range egressAdapterErrors {
				log.Printf("Invocating %s failed: %s",
					reflect.TypeOf(egressAdapterError.EgressAdapter).Elem().Name(),
					egressAdapterError.Error,
				)
			}

			restResponse := contract.RestResponseBody{
				Message: fmt.Sprintf("Successfully ran %s.", handler.ID),
			}
			return c.Status(fiber.StatusOK).JSON(restResponse)
		})
	default:
		err = errors.New("Handler has no matching action")
	}
	return err
}
