package services

import (
	"fmt"

	"github.com/wahlly/Digiwallet-demo/models"
	"github.com/wahlly/Digiwallet-demo/modules/paystack"
)

type WalletService struct {
	us *UserService
}

func NewWalletService(us *UserService) *WalletService {
	return &WalletService{us: us}
}

func (ws *WalletService) GetUserWallet(id uint) (*models.Wallet, error) {
	var user models.User
	err := ws.us.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user.Wallet, nil
}

func (ws *WalletService) InitializeWalletDeposit(payload *paystack.InitTxnReqBody) (map[string]any, error) {
	paystackClient := paystack.NewPaystackClient()
	res, err := paystackClient.InitiateTransaction(payload.Amount, payload.Email)
	if err != nil {
		return nil, err
	}
fmt.Println("resp: ", res)
	return res, nil
}