package models

import "gorm.io/gorm"

type Banner struct {
	gorm.Model
	ImageURL string
	Title    string
	Subtitle string
	Link     string
}
