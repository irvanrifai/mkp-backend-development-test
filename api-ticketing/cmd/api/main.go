package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/database"
	authHandler "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/delivery/http/auth"
	eventHandler "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/delivery/http/event"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/middleware"
	eventRepository "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/repository/event"
	userRepository "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/repository/user"
	eventUsecase "github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/usecase/event"
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

	eventRepo := eventRepository.NewEventRepository(db)
	eventUC := eventUsecase.NewEventUsecase(eventRepo)
	eventHandler := eventHandler.NewEventHandler(eventUC)

	app := fiber.New()
	api := app.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Event Routes
	events := api.Group("/events")

	// Public: List & Detail
	events.Get("/", eventHandler.List)
	events.Get("/:id", eventHandler.Detail)

	// Protected: Create, Update, Delete
	events.Post("/", middleware.Protected, eventHandler.Create)
	events.Put("/:id", middleware.Protected, eventHandler.Update)
	events.Delete("/:id", middleware.Protected, eventHandler.Delete)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
