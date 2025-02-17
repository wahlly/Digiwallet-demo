package utils

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/wahlly/Digiwallet-demo/models"
	"gorm.io/gorm"
)


func GenerateWalletAddress(db *gorm.DB) (string, error) {
	address := ""
	for i := 0; i < 10; i++ {
		address += strconv.Itoa(rand.Intn(10))
	}

	var user *models.User
	err := db.Where("wallet->>'address' = ?", address).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GenerateWalletAddress(db)
		} else{
			return "", err
		}
	}

	return address, nil
}