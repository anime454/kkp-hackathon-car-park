package main

import (
	"log"

	httpadapter "be/internal/adapter/in/http"
	"be/internal/adapter/out/persistence/postgres"
	"be/internal/config"
	"be/internal/core/domain"
	"be/internal/core/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	db, err := postgres.Connect(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Gorm().AutoMigrate(&domain.CarPark{}); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	healthSvc := service.NewHealthService(db)
	carParkRepo := postgres.NewCarParkRepository(db.Gorm())
	carParkSvc := service.NewCarParkService(carParkRepo)

	httpadapter.RegisterRoutes(app, healthSvc, carParkSvc)

	log.Fatal(app.Listen(cfg.HTTP.Addr()))
}
