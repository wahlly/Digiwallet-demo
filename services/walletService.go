package services

import (
	"errors"

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

	if !res.Status {
		return nil, errors.New(res.Message)
	}

	return map[string]any{"data": res}, nil
}