package models

import "gorm.io/gorm"

type KonuluKonum struct {
	gorm.Model
	URI             string
	ImageURI        string
	Loc             string
	AuthorName      string
	Description     string
	UnlockedCounter int
}
