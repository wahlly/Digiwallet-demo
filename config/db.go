package config

import (
	"github.com/wahlly/Digiwallet-demo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB *gorm.DB

func ConnectDB () {
	db, err := gorm.Open(postgres.Open("postgres://wahlly:369852@localhost:5432/digiWal"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})
	DB = db

}