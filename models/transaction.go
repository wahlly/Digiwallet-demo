package models

import "gorm.io/gorm"



type Transaction struct {
	gorm.Model
	userId
}