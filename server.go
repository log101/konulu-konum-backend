package main

import (
	"fmt"
	"io"
	"log"
	DB "log101/konulu-konum-backend/db"
	"log101/konulu-konum-backend/models"

	"os"

	"github.com/dchest/uniuri"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/h2non/bimg"
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
	db := DB.GetDB()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Sample")
	})

	app.Post("/api/location", func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			// Get string input
			if token := form.Value["token"]; len(token) > 0 {
				fmt.Println(token[0])
			}

			// Get image
			files := form.File["document"]
			if len(files) != 1 {
				fmt.Println("bad request")
			}

			file := form.File["document"][0]
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

			newFile, err := file.Open()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			defer newFile.Close()

			data, err := io.ReadAll(newFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			// Compress image
			newImage, err := bimg.NewImage(data).Convert(bimg.WEBP)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			// Generate public uri for the image
			chars := uniuri.StdChars[26:52]
			randomUri := uniuri.NewLenChars(10, chars)
			imageUri := fmt.Sprintf("%s-%s-%s", randomUri[0:3], randomUri[3:7], randomUri[7:])

			db.Create(&models.KonuluKonum{URI: imageUri, Image: newImage, Coordinates: "sample", AuthorName: "sample", Description: "sample", UnlockedCounter: 0})

			return c.SendStatus(fiber.StatusOK)
		}

		return c.SendStatus(fiber.StatusBadRequest)
	})

	app.Get("/api/location/:locationUri", func(c *fiber.Ctx) error {
		uri := c.Params("locationUri")
		if len(uri) == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var konuluKonum models.KonuluKonum
		rows := db.Where("URI = ?", uri).First(&konuluKonum)
		if rows.Error != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendString(konuluKonum.URI)
	})

	app.Listen(":3000")

}
