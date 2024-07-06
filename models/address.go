package models

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	CustomerID    uint   `json:"customer_id"`
	ReceiverName  string `json:"receiver_name" form:"receiver_name"`
	PhoneNumber   string `json:"phone_number" form:"phone_number"`
	HouseNumber   string `json:"house_number" form:"house_number"`
	ResidenceName string `json:"residence_name" form:"residence_name"`
	AddressNotes  string `json:"address_notes" form:"address_notes"`
	StreetName    string `json:"street_name" form:"street_name"`
	District      string `json:"district" form:"district"`
	SubDistrict   string `json:"sub_district" form:"sub_district"`
	City          string `json:"city" gorm:"default:'Bandung'"`
	Area          string `json:"area" gorm:"default:'Bojongsoang'"`
}

func (address *Address) BeforeCreate(tx *gorm.DB) (err error) {
	if address.City == "" {
		address.City = "Bandung"
	}
	if address.Area == "" {
		address.Area = "Bojongsoang"
	}
	return
}
