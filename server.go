package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Sample")
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			if token := form.Value["token"]; len(token) > 0 {
				fmt.Println(token[0])
			}

			files := form.File["documents"]

			for _, file := range files {
				fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			}

			return err
		}

		return nil
	})

	app.Listen(":3000")
}
