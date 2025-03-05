package models

import "gorm.io/gorm"



type Transaction struct {
	gorm.Model
	Amount	uint	`json:"amount"`
	Type 		string `json:"type"`
	Ref		string `json:"ref"`
	UserId	int	`json:"userId"`
	User		User
}