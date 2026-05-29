package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/database"
	authHandler "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/delivery/http/auth"
	scheduleHandler "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/delivery/http/schedule"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/middleware"
	scheduleRepository "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/repository/schedule"
	userRepository "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/repository/user"
	scheduleUsecase "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/usecase/schedule"
	userUsecase "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/usecase/user"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()
	db := database.DB

	userRepo := userRepository.NewUserRepository(db)
	userUC := userUsecase.NewUserUsecase(userRepo, os.Getenv("JWT_SECRET"))
	authHandler := authHandler.NewAuthHandler(userUC)

	scheduleRepo := scheduleRepository.NewScheduleRepository(db)
	scheduleUC := scheduleUsecase.NewScheduleUsecase(scheduleRepo)
	scheduleHandler := scheduleHandler.NewScheduleHandler(scheduleUC)

	app := fiber.New()
	api := app.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Schedule Routes
	schedules := api.Group("/schedules")
	schedules.Get("/", middleware.Protected, scheduleHandler.List)
	schedules.Post("/", middleware.Protected, scheduleHandler.Create)
	schedules.Get("/:id", middleware.Protected, scheduleHandler.Detail)
	schedules.Put("/:id", middleware.Protected, scheduleHandler.Update)
	schedules.Delete("/:id", middleware.Protected, scheduleHandler.Delete)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
