package services

import (
	"context"
	"errors"
	"strings"

	"github.com/wahlly/Digiwallet-demo/models"
	"github.com/wahlly/Digiwallet-demo/modules/paystack"
	"github.com/wahlly/Digiwallet-demo/utils"
	"gorm.io/gorm"
)

type WalletService struct {
	us *UserService
	DB *gorm.DB
}

func NewWalletService(us *UserService, db *gorm.DB) *WalletService {
	return &WalletService{us: us, DB: db}
}

func (ws *WalletService) GetUserWallet(id uint) (*models.Wallet, error) {
	var user models.User
	err := ws.us.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user.Wallet, nil
}

func (ws *WalletService) GetWalletByAddress(address string) (*models.Wallet, error) {
	var user models.User
	err := ws.us.db.Where("wallet->>'address' = ?", address).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user.Wallet, nil
}

func (ws *WalletService) InitializeWalletDeposit(payload *paystack.InitTxnReqBody, userId uint) (*paystack.PaystackInitTxnRes, error) {
	if payload.Amount < 100 {
		return nil, errors.New("minimum deposit allowed is 100")
	}
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
		User_id: uint(userId),
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

	var txn models.Transaction
	err = ws.us.db.Where("reference = ? AND user_id = ?", reference, uid).First(&txn).Error
	if err != nil {
		return nil, err
	}

	if txn.Status != "pending" {
		return nil, errors.New("transaction has already been settled")
	}

	res, err = paystackClient.VerifyTransaction(uid, reference)
	if err != nil {
		return nil, err
	}

	if !res.Status {
		return nil, errors.New(res.Message) 
	}

	if strings.ToLower(res.Data["gateway_response"].(string)) != "successful" {
		return nil, errors.New(res.Data["gateway_response"].(string))
	}

	txn.Status = "success"
	ws.us.db.Save(&txn)

	var user models.User
	err = ws.us.db.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}

	user.Wallet.Balance += int64(txn.Amount)
	ws.us.db.Save(&user)

	return map[string]any{"amount": txn.Amount/100}, nil
}

func (ws *WalletService) TransferToWalletAddress(ctx context.Context, dbTxn *gorm.DB, uid uint, amount int64, recipientAddress string) error {
	var dbTxnWithCtx = dbTxn.WithContext(ctx)	//db transaction session with context
	var user models.User
	err := dbTxnWithCtx.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return err
	}

	if user.Wallet.Balance < amount * 100 {
		return errors.New("insufficient account balance")
	}

	var recipient models.User
	err = dbTxnWithCtx.Where("wallet->>'address' = ?", recipientAddress).First(&recipient).Error
	if err != nil {
		return err
	}

	user.Wallet.Balance -= amount * 100
	err = dbTxnWithCtx.Save(&user).Error
	if err != nil{
		return err
	}

	recipient.Wallet.Balance += amount * 100
	err = dbTxnWithCtx.Save(&recipient).Error
	if err != nil{
		return err
	}

	//log transaction for credit and debit
	debit := &models.Transaction{
		Amount: amount * 100,
		User_id: uint(user.ID),
		Type: "debit",
		Reference: strings.ToUpper(utils.AlphaNumeric(10, "alphaNumeric")),
		Status: "success",
	}
	err = dbTxnWithCtx.Create(debit).Error
	if err != nil {
		return err
	}
	credit := &models.Transaction{
		Amount: amount * 100,
		User_id: uint(recipient.ID),
		Type: "credit",
		Reference: strings.ToUpper(utils.AlphaNumeric(10, "alphaNumeric")),
		Status: "success",
	}
	dbTxnWithCtx.Create(credit)

	return nil
}

//transfer to banks