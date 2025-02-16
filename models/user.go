package models

import "gorm.io/gorm"

type wallet struct {
	Address	string	`json:"address" gorm:"unique"`
	Balance	int64		`json:"balance" gorm:"not null;default:0"`
}

type User struct {
	gorm.Model
	FirstName	string	`json:"firstName"`
	LastName	string	`json:"lastName"`
	Phone		uint8		`json:"phone" gorm:"unique"`
	Email		string	`json:"email" gorm:"unique"`
	UserName	string	`json:"userName" gorm:"unique"`
	Password	string	`json:"password" gorm:"not null"`
	Wallet 	wallet	`json:"wallet" gorm:"type:json"`
}