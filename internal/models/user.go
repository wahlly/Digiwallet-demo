package models

import (
	// "database/sql/driver"
	// "encoding/json"
	// "errors"

	"gorm.io/gorm"
)

type Wallet struct {
	Address	string	`json:"address" gorm:"unique"`
	Balance	int64		`json:"balance" gorm:"not null;default:0"`
}

// func (w Wallet) Value() (driver.Value, error) {
// 	return json.Marshal(w)
// }

// func (w *Wallet) Scan(value interface{}) error {
// 	if value == nil {
// 		return nil
// 	}

// 	bytes, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("failed to scan Wallet JSON")
// 	}

// 	return json.Unmarshal(bytes, w)
// }

type User struct {
	gorm.Model
	FirstName	string	`json:"firstName"`
	LastName	string	`json:"lastName"`
	Phone		string	`json:"phone" gorm:"unique"`
	Email		string	`json:"email" gorm:"unique"`
	UserName	string	`json:"userName" gorm:"unique"`
	Password	string	`json:"password" gorm:"not null"`
	Wallet 	Wallet	`json:"wallet" gorm:"type:jsonb;serializer:json"`
}