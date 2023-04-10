package rest

import (
	"github.com/gofiber/fiber/v2"
)

type TestResponse struct {
	Message string `json:"message"`
}

func getTestHandler(c *fiber.Ctx) error {
	response := TestResponse{
		Message: "GET request received on /test",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func postTestHandler(c *fiber.Ctx) error {
	response := TestResponse{
		Message: "POST request received on /test",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
