package rest

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
)

// must implement core.ServerAdapter
type RestServerAdapter struct {
	addr   string
	server *fiber.App
}

func NewRestServerAdapter() RestServerAdapter {
    app := fiber.New()
    app.Use(logger.New())
    
	return RestServerAdapter{addr: ":8080", server: app}
}

func (rsa RestServerAdapter) Serve() error {
	return rsa.server.Listen(rsa.addr)
}

type RestResponse struct {
	Message string `json:"message"`
}

func (rsa RestServerAdapter) Register(handler *core.Handler) error {
	var err error = nil
	switch handler.Action {
	case "read":
		rsa.server.Get("/"+handler.Name, func(c *fiber.Ctx) error {
			for _, task := range handler.TaskSets {
				_, err := task.ServiceAdapter.Call()
				if err != nil {
					return err
				}
			}
			return c.Status(fiber.StatusOK).JSON(RestResponse{Message: fmt.Sprintf("Successfully ran %s", handler.ID)})
		})
	case "write":
		rsa.server.Post("/"+handler.Name, func(c *fiber.Ctx) error {
			for _, task := range handler.TaskSets {
				_, err := task.ServiceAdapter.Call()
				if err != nil {
					return err
				}
			}
			return c.Status(fiber.StatusOK).JSON(RestResponse{Message: fmt.Sprintf("Successfully ran %s", handler.ID)})
		})
	default:
		err = errors.New("Handler has no matching action")
	}
	return err
}
