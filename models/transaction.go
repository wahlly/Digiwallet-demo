package models

import "gorm.io/gorm"



type Transaction struct {
	gorm.Model
	Amount		uint		`json:"amount"`
	Type 			string	`json:"type"`
	Reference		string	`json:"reference"`
	Status		string	`json:"status"`
	User_id		int		`json:"user_id"`
	User			User
}