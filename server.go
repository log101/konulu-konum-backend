package main

import (
	"fmt"
	"io"
	"os"
	"strings"

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
			if token := form.Value["token"]; len(token) > 0 {
				fmt.Println(token[0])
			}

			files := form.File["documents"]

			for _, file := range files {
				fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

				// Save the files to disk:
				//				if err := c.SaveFile(file, fmt.Sprintf("./public/%s", file.Filename)); err != nil {
				//					return err
				//				}

				/*
					buffer, err := bimg.Read(file.Filename)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				*/

				newFile, err := file.Open()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				defer newFile.Close()

				data, err := io.ReadAll(newFile)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				newImage, err := bimg.NewImage(data).Convert(bimg.WEBP)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				imageName := strings.Split(file.Filename, ".")[0]

				bimg.Write(fmt.Sprintf("./public/%s.webp", imageName), newImage)
			}

			return err
		}

		return nil
	})

	app.Listen(":3000")

}
