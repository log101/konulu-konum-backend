package main

import (
	"log"
	DB "log101/konulu-konum-backend/db"
	"log101/konulu-konum-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// initialize db
	DB.InitDB()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4321",
	}))
	app.Static("/images", "./public")

	app.Post("/api/location", handlers.KonuluKonumCreate)

	app.Get("/api/location/:locationUri", handlers.KonuluKonumGet)

	app.Patch("/api/location/increment/:locationUri", handlers.KonuluKonumUpdateCounter)

	app.Listen(":3456")

}
