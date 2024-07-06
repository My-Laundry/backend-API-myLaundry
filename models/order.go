package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID uint    `json:"customer_id"`
	CourierID  uint    `json:"courier_id"`
	AdminID    uint    `json:"admin_id"`
	ServiceID  uint    `json:"service_id"`
	AddressID  uint    `json:"address_id"`
	Weight     float64 `json:"weight"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"` // "waiting for pickup", "in process", "out for delivery", "completed"
	Address    Address `json:"address" gorm:"foreignKey:AddressID"`
	Customer   User    `json:"customer" gorm:"foreignKey:CustomerID"`
	Courier    User    `json:"courier" gorm:"foreignKey:CourierID"`
	Admin      User    `json:"admin" gorm:"foreignKey:AdminID"`
	Service    Service `json:"service" gorm:"foreignKey:ServiceID"`
}
