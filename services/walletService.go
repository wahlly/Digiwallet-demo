package services

import (
	"github.com/wahlly/Digiwallet-demo/models"
)

type WalletService struct {
	us *UserService
}

func NewWalletService(us *UserService) *WalletService {
	return &WalletService{us: us}
}

func (ws *WalletService) GetUserWallet(email string) (*models.Wallet, error) {
	var user models.User
	err := ws.us.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user.Wallet, nil
}