package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     int64   `json:"price"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
