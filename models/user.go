// package models

// import (
// 	"database/sql/driver"
// 	"encoding/json"
// 	"errors"

// 	"gorm.io/gorm"
// )

// type User struct {
// 	gorm.Model
// 	Username  string    `json:"username"`
// 	Email     string    `json:"email" gorm:"unique"`
// 	Password  string    `json:"password"`
// 	Role      string    `json:"role"` // "customer", "admin", "courier"
// 	Addresses Addresses `gorm:"-"`    // Tidak ada tag gorm untuk Addresses
// }

// type Addresses []Address

// // Implement Valuer interface
// func (a Addresses) Value() (driver.Value, error) {
// 	return json.Marshal(a)
// }

// // Implement Scanner interface
// func (a *Addresses) Scan(value interface{}) error {
// 	if data, ok := value.([]byte); ok {
// 		return json.Unmarshal(data, a)
// 	}
// 	return errors.New("failed to scan addresses")
// }
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
