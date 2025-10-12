package config

import (
	"github.com/wahlly/Digiwallet-demo/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB () *gorm.DB {
	db, err := gorm.Open(
		postgres.Open("postgres://postgres:admin@localhost:5432/digiWal"),
		&gorm.Config{
			SkipDefaultTransaction: true,	//undo db write transaction by default
		},
	)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
	)

	return db
}