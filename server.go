package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/h2non/bimg"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Sample")
	})

	app.Static("/images", "./public")

	app.Post("/upload", func(c *fiber.Ctx) error {
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

			imageName := strings.Split(file.Filename, ".")[0]

			// Save image
			bimg.Write(fmt.Sprintf("./public/%s.webp", imageName), newImage)

			// Generate public uri for the image
			chars := uniuri.StdChars[26:52]
			randomUri := uniuri.NewLenChars(10, chars)
			imageUri := fmt.Sprintf("%s-%s-%s", randomUri[0:3], randomUri[3:7], randomUri[7:])

			fmt.Println(imageUri)

			return c.SendStatus(fiber.StatusOK)
		}

		return c.SendStatus(fiber.StatusBadRequest)
	})

	app.Listen(":3000")

}
