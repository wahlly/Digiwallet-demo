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
	amount := payload.Amount + paystackClient.CalculateProcessingFee(payload.Amount)
	res, err = paystackClient.InitiateTransaction(amount, payload.Email, ref)
	if err != nil {
		return res, err
	}

	if !res.Status {
		return nil, errors.New(res.Message)
	}

	txn := &models.Transaction{
		Amount: payload.Amount * 100,
		User_id: int(userId),
		Type: "credit",
		Reference: ref,
	}
	err = ws.us.db.Create(txn).Error
	if err != nil {
		return res, err
	}

	// return map[string]any{"data": res}, nil
	return res, nil
}

func (ws *WalletService) VerifyWalletDeposit(uid uint, reference string) (map[string]any, error) {
	paystackClient := paystack.NewPaystackClient()
	res := &paystack.PaystackVerifyTxnRes{}
	var err error
	res, err = paystackClient.VerifyTransaction(uid, reference)
	if err != nil {
		return nil, err
	}

	if !res.Status {
		return nil, errors.New(res.Message) 
	}

	var txn models.Transaction
	err = ws.us.db.Where("reference = ? AND user_id = ?", reference, uid).First(&txn).Error
	if err != nil {
		return nil, err
	}

	txn.Status = "success"
	ws.us.db.Save(&txn)

	var user models.User
	err = ws.us.db.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}

	user.Wallet.Balance += int64(txn.Amount)

	return map[string]any{"amount": txn.Amount}, nil
}