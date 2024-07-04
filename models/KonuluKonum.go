package models

import (
	"gorm.io/gorm"
)

type KonuluKonum struct {
	gorm.Model
	URI             string
	Image           []byte `gorm:"type:BLOB"`
	Coordinates     string
	AuthorName      string
	Description     string
	UnlockedCounter int
}
