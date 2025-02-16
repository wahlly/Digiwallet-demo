package services

import (
	"errors"

	"github.com/wahlly/Digiwallet-demo/models"
	"github.com/wahlly/Digiwallet-demo/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{
	db *gorm.DB
}

func (us *UserService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil
	}

	user.Password = string(hashedPassword)
	err = us.db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) LoginUser(email, password string) (id uint, token string, err error) {
	var user models.User
	err = us.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, "", errors.New("email or password is incorrect")
	}

	token, err = utils.GenerateJWTtoken(user)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}