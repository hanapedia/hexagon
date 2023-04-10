package rest

import (
	"github.com/gofiber/fiber/v2"
)

func setupRouter() *fiber.App {
	app := fiber.New()

	app.Get("/test", getTestHandler)
	app.Post("/test", postTestHandler)

	return app
}
