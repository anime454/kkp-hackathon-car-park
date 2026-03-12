package http

import (
	"context"

	portin "be/internal/core/port/in"
	"github.com/gofiber/fiber/v2"
)

type createCarParkRequest struct {
	Name string `json:"name"`
}

func RegisterRoutes(app *fiber.App, health portin.HealthUseCase, carParks portin.CarParkUseCase) {
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := health.Check(context.Background()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
			})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/carparks", func(c *fiber.Ctx) error {
		var req createCarParkRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		}
		if req.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
		}

		cp, err := carParks.Create(context.Background(), req.Name)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create"})
		}
		return c.Status(fiber.StatusCreated).JSON(cp)
	})

	app.Get("/carparks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		cp, err := carParks.Get(context.Background(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}
		return c.JSON(cp)
	})
}
