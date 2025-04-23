package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID       uint
	FirstName    string
	LastName     string
	Address      string
	City         string
	Email        string
	Phone        string
	Payment      string
	OrderNotes   string
	Status       string
	OrderDetails []OrderDetail
	Total        int64 `gorm:"-"`
}

type OrderDetail struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Product   Product
	Quantity  int
	Price     int64
}
