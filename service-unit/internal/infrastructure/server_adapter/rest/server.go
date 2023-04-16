package rest

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/contract"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/server_adapter/common"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
	"github.com/hanapedia/the-bench/service-unit/pkg/utils"
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
			invocationAdapterErrors := common.TaskSetHandler(handler.TaskSets)
			for _, invocationAdapterError := range invocationAdapterErrors {
				log.Printf("Invocating %s failed: %s",
					reflect.TypeOf(invocationAdapterError.InvocationAdapter).Elem().Name(),
					invocationAdapterError.Error,
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
			invocationAdapterErrors := common.TaskSetHandler(handler.TaskSets)
			for _, invocationAdapterError := range invocationAdapterErrors {
				log.Printf("Invocating %s failed: %s",
					reflect.TypeOf(invocationAdapterError.InvocationAdapter).Elem().Name(),
					invocationAdapterError.Error,
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
