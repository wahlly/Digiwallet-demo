package models

import "gorm.io/gorm"



type Transaction struct {
	gorm.Model
	Amount		uint		`json:"amount"`
	Type 			string	`json:"type"`
	Reference		string	`json:"reference"`
	Status		string	`json:"status" gorm:"default:pending"`
	User_id		int		`json:"user_id"`
	User			User
}