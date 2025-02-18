package utils

import (
	"math/rand"
	"strconv"

	"github.com/wahlly/Digiwallet-demo/models"
	"gorm.io/gorm"
)


func GenerateWalletAddress(db *gorm.DB) (string, error) {
	for {
		var address string
		for i := 0; i < 10; i++ {
			address += strconv.Itoa(rand.Intn(10))
		}

		var count int64
		err := db.Model(&models.User{}).Where("wallet->>'address' = ?", address).Count(&count).Error
		if err != nil {
			return "", err
		}

		if count == 0 {
			return address, nil
		}
	}
}