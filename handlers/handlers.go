package handlers

import (
	"fmt"
	"io"
	"os"
	"strings"

	DB "log101/konulu-konum-backend/db"
	"log101/konulu-konum-backend/models"

	"github.com/dchest/uniuri"
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/bimg"
	"gorm.io/gorm"
)

func KonuluKonumCreate(c *fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		// Get form values
		author := form.Value["author"][0]
		description := form.Value["description"][0]
		geolocation := fmt.Sprintf("[%s]", form.Value["geolocation"][0])

		file := form.File["selected-photo"][0]
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
		imagePath := fmt.Sprintf("./public/%s.webp", imageName)
		imageURL := fmt.Sprintf("%s.webp", imageName)
		err = bimg.Write(imagePath, newImage)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		// Generate public uri for the image
		chars := uniuri.StdChars[26:52]
		randomUri := uniuri.NewLenChars(10, chars)
		imageUri := fmt.Sprintf("%s-%s-%s", randomUri[0:3], randomUri[3:7], randomUri[7:])

		db := DB.GetDB()
		db.Create(&models.KonuluKonum{URI: imageUri, ImageURL: imageURL, Coordinates: geolocation, AuthorName: author, Description: description, UnlockedCounter: 0})

		return c.JSON(fiber.Map{
			"url": imageUri,
		})
	}

	return c.SendStatus(fiber.StatusBadRequest)
}

func KonuluKonumGet(c *fiber.Ctx) error {
	uri := c.Params("locationUri")
	if len(uri) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var konuluKonum models.KonuluKonum
	db := DB.GetDB()
	rows := db.Where("URI = ?", uri).First(&konuluKonum)
	if rows.Error != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"url":              konuluKonum.URI,
		"blob_url":         konuluKonum.ImageURL,
		"loc":              konuluKonum.Coordinates,
		"author":           konuluKonum.AuthorName,
		"description":      konuluKonum.Description,
		"unlocked_counter": konuluKonum.UnlockedCounter,
	})
}

func KonuluKonumUpdateCounter(c *fiber.Ctx) error {
	uri := c.Params("locationUri")
	if len(uri) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var konuluKonum models.KonuluKonum
	db := DB.GetDB()
	rows := db.Where("URI = ?", uri).First(&konuluKonum)
	if rows.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	rows = db.Model(&konuluKonum).Where("uri = ?", uri).UpdateColumn("unlocked_counter", gorm.Expr("unlocked_counter + 1"))
	if rows.Error != nil {
		c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"counter": konuluKonum.ID,
	})
}
