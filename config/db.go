package config

import (
	"github.com/wahlly/Digiwallet-demo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB () *gorm.DB {
	db, err := gorm.Open(postgres.Open("postgres://postgres:admin@localhost:5432/digiWal"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	return db
}