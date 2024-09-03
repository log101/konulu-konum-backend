package handlers

import (
	"fmt"
	"io"
	"os"
	"strconv"

	DB "log101/konulu-konum-backend/db"
	"log101/konulu-konum-backend/models"

	"github.com/dchest/uniuri"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"gorm.io/gorm"
)

func KonuluKonumCreate(c *fiber.Ctx) error {
	clientURL := os.Getenv("CLIENT_URL")
	if form, err := c.MultipartForm(); err == nil {
		// Get form values
		author := form.Value["author"][0]
		description := form.Value["description"][0]
		radius := form.Value["geolocation-radius"][0]
		radiusInt, err := strconv.Atoi(radius)
		if err != nil {
			radiusInt = 50
		}

		// Geolocation is stored as JSON array string
		geolocation := fmt.Sprintf("[%s]", form.Value["geolocation"][0])
		file := form.File["selected-photo"]
		if len(file) != 1 {
			fmt.Fprintln(os.Stderr, err)
			redirectUrl := fmt.Sprintf("%s?error=%s", clientURL, "true")
			return c.Redirect(redirectUrl)
		}

		newFile, err := file[0].Open()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			redirectUrl := fmt.Sprintf("%s?error=%s", clientURL, "true")
			return c.Redirect(redirectUrl)
		}
		defer newFile.Close()

		// Read image file
		data, err := io.ReadAll(newFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			redirectUrl := fmt.Sprintf("%s?error=%s", clientURL, "true")
			return c.Redirect(redirectUrl)
		}

		// Compress image file and convert to webp
		newImage, err := bimg.NewImage(data).Convert(bimg.WEBP)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			redirectUrl := fmt.Sprintf("%s?error=%s", clientURL, "true")
			return c.Redirect(redirectUrl)
		}

		// Save image file in public folder
		imageName := uuid.New()
		imagePath := fmt.Sprintf("./public/%s.webp", imageName)
		imageNameWithExtension := fmt.Sprintf("%s.webp", imageName)
		err = bimg.Write(imagePath, newImage)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			redirectUrl := fmt.Sprintf("%s?error=%s", clientURL, "true")
			return c.Redirect(redirectUrl)
		}

		// Generate public uri for the image this will be the
		// id for the konulu konum
		chars := uniuri.StdChars[26:52]
		randomUri := uniuri.NewLenChars(10, chars)
		imageUri := fmt.Sprintf("%s-%s-%s", randomUri[0:3], randomUri[3:7], randomUri[7:])

		// Write to DB
		db := DB.GetDB()
		db.Create(&models.KonuluKonum{URI: imageUri, ImageURL: imageNameWithExtension, Coordinates: geolocation, AuthorName: author, Description: description, UnlockedCounter: 0, Radius: radiusInt})

		// Return URL
		redirectURL := fmt.Sprintf("%s/x?id=%s", clientURL, imageUri)
		return c.Redirect(redirectURL)
	}

	redirectUrl := fmt.Sprintf("%s?error=%s", clientURL, "true")
	return c.Redirect(redirectUrl)
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
		"radius":           konuluKonum.Radius,
		"unlocked_counter": konuluKonum.UnlockedCounter,
		"created_at":       konuluKonum.CreatedAt,
	})
}

func KonuluKonumCounterUpdate(c *fiber.Ctx) error {
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
