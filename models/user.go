package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique" json:"email"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"password"`
	Address   string    `gorm:"not null" json:"address"`
	Birthdate time.Time `gorm:"type:date" json:"birthdate"`
	Role      string    `gorm:"type:varchar(20);default:'customer'" json:"role"`
}

type Info struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Address   string `json:"address"`
	Birthdate string `json:"birthdate"`
}
