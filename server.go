package main

import (
	"log"
	DB "log101/konulu-konum-backend/db"
	"log101/konulu-konum-backend/handlers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Create public folder for images
	err = os.Mkdir("public", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Initialize db
	DB.InitDB()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4321",
	}))

	// Serve static images
	app.Static("/images", "./public")

	// Create konulu konum
	app.Post("/api/location", handlers.KonuluKonumCreate)

	// Get konulu konum
	app.Get("/api/location/:locationUri", handlers.KonuluKonumGet)

	// Update 'seen' counter of konulu konum
	// This is shown at the bottom of the web page
	app.Patch("/api/location/increment/:locationUri", handlers.KonuluKonumCounterUpdate)

	app.Listen(":3456")
}
