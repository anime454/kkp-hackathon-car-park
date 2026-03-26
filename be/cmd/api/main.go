package main

import (
	"context"
	"log"

	httpadapter "be/internal/adapter/in/http"
	"be/internal/adapter/out/persistence/postgres"
	"be/internal/config"
	"be/internal/core/domain"
	"be/internal/core/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.Load()

	db, err := postgres.Connect(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Gorm().AutoMigrate(&domain.CarPark{}, &domain.Slot{}); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(cors.New())

	healthSvc := service.NewHealthService(db)
	carParkRepo := postgres.NewCarParkRepository(db.Gorm())
	slotRepo := postgres.NewSlotRepository(db.Gorm())
	carParkSvc := service.NewCarParkService(carParkRepo, slotRepo, cfg.Admin)
	if err := carParkSvc.EnsureSeedData(context.Background()); err != nil {
		log.Fatal(err)
	}

	httpadapter.RegisterRoutes(app, healthSvc, carParkSvc)

	log.Fatal(app.Listen(cfg.HTTP.Addr()))
}
