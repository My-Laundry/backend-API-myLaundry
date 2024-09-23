package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	Role      string    `json:"role"` // "customer", "admin", "courier"
	Addresses []Address `gorm:"foreignkey:CustomerID"`
}
