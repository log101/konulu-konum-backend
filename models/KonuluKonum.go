package models

import (
	"gorm.io/gorm"
)

type KonuluKonum struct {
	gorm.Model
	URI             string
	ImageURL        string
	Coordinates     string
	AuthorName      string
	Description     string
	Radius          int
	UnlockedCounter int
}
