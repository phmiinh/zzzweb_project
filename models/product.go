package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(255);not null"`
	ImageURL    string  `gorm:"type:varchar(255);not null"`
	Price       float64 `gorm:"not null"`
	OldPrice    float64 `gorm:""`
	Stock       int     `gorm:"not null;default:0"`
	Category    string  `json:"category"`
	Description string  `gorm:"type:text"`
}
