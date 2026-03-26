package http

import (
	"context"
	"strings"
	"time"

	portin "be/internal/core/port/in"
	"github.com/gofiber/fiber/v2"
)

type createCarParkRequest struct {
	Name string `json:"name"`
}

type adminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type updateSlotStatusRequest struct {
	Status string `json:"status"`
}

func bearerToken(c *fiber.Ctx) string {
	auth := strings.TrimSpace(c.Get("Authorization"))
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
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

	app.Get("/kiosk/dashboard", func(c *fiber.Ctx) error {
		dashboard, err := carParks.GetDashboard(context.Background())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch dashboard"})
		}
		return c.JSON(dashboard)
	})

	app.Get("/kiosk/info", func(c *fiber.Ctx) error {
		dashboard, err := carParks.GetDashboard(context.Background())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch info"})
		}
		return c.JSON(fiber.Map{
			"message":       "Welcome to KKP Smart Parking",
			"updatedAt":     time.Now().UTC(),
			"vipRate":       "20 first hour, +10 each additional hour",
			"normalRate":    "10 first hour, +5 each additional hour",
			"availableSlot": dashboard.Free,
			"occupiedSlot":  dashboard.Parked,
		})
	})

	app.Get("/parking/fee", func(c *fiber.Ctx) error {
		slotID := c.Query("slotId")
		if slotID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "slotId is required"})
		}
		exitAt := time.Now().UTC()
		if raw := strings.TrimSpace(c.Query("exitAt")); raw != "" {
			parsed, err := time.Parse(time.RFC3339, raw)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "exitAt must be RFC3339"})
			}
			exitAt = parsed
		}

		fee, err := carParks.CalculateFee(context.Background(), slotID, exitAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fee)
	})

	app.Post("/management/login", func(c *fiber.Ctx) error {
		var req adminLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		}
		token, err := carParks.AdminLogin(context.Background(), req.Username, req.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return c.JSON(fiber.Map{"token": token})
	})

	management := app.Group("/management", func(c *fiber.Ctx) error {
		token := bearerToken(c)
		if err := carParks.Authorize(context.Background(), token); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Next()
	})

	management.Get("/slots", func(c *fiber.Ctx) error {
		slots, err := carParks.ListSlots(context.Background())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch slots"})
		}
		return c.JSON(fiber.Map{"slots": slots})
	})

	management.Patch("/slots/:id/status", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req updateSlotStatusRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		}
		updated, err := carParks.UpdateSlotStatus(context.Background(), id, req.Status)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(updated)
	})
}
