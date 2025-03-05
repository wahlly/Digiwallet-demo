package services

import (
	"errors"
	"strings"

	"github.com/wahlly/Digiwallet-demo/models"
	"github.com/wahlly/Digiwallet-demo/modules/paystack"
	"github.com/wahlly/Digiwallet-demo/utils"
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

func (ws *WalletService) InitializeWalletDeposit(payload *paystack.InitTxnReqBody, userId uint) (*paystack.PaystackInitTxnRes, error) {
	ref := strings.ToUpper(utils.AlphaNumeric(10, "alphaNumeric"))
	res := &paystack.PaystackInitTxnRes{}
	var err error
	paystackClient := paystack.NewPaystackClient()
	res, err = paystackClient.InitiateTransaction(payload.Amount, payload.Email, ref)
	if err != nil {
		return res, err
	}

	if !res.Status {
		return res, errors.New(res.Message)
	}

	txn := &models.Transaction{
		Amount: payload.Amount * 100,
		UserId: int(userId),
		Type: "credit",
		Ref: ref,
	}
	err = ws.us.db.Create(txn).Error
	if err != nil {
		return res, err
	}

	// return map[string]any{"data": res}, nil
	return res, nil
}