package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/modules/paystack"
	"github.com/wahlly/Digiwallet-demo/services"
	"github.com/wahlly/Digiwallet-demo/utils"
)


type WalletController struct {
	WalletService *services.WalletService
}

func (wc *WalletController) GetUserWallet(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, sign in again"})
		c.Abort()
		return
	}

	wallet, err := wc.WalletService.GetUserWallet(id.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "wallet": wallet})
}

func (wc *WalletController) GetWalletByAddress(c *gin.Context) {
	address := c.Param("address")
	wallet, err := wc.WalletService.GetWalletByAddress(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "wallet": wallet})
}

func (wc *WalletController) InitializeWalletDeposit(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, sign in again"})
		c.Abort()
		return
	}

	var payload paystack.InitTxnReqBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	res, err := wc.WalletService.InitializeWalletDeposit(&payload, id.(uint))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	r := map[string]any{
		"uid": id,
		"payment_res": res,
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": r})
}

func (wc *WalletController) VerifyWalletDeposit(c *gin.Context) {
	reference := c.Param("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid reference"})
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, sign in again"})
		return
	}

	res, err := wc.WalletService.VerifyWalletDeposit(id.(uint), reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wallet deposit successfull", "data": res})
}

func (wc *WalletController) TransferToWalletAddress(c *gin.Context) {
	type reqBody struct{
		Amount int64 `json:"amount" binding:"required,min=100"`
		Recipient string	`json:"recipient" binding:"required"`
	}
	var payload reqBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, sign in again"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx := wc.WalletService.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to start transaction"})
	}
	defer func () {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
	}()

	err := wc.WalletService.TransferToWalletAddress(ctx, tx, id.(uint), payload.Amount, payload.Recipient)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to complete transaction"})
		return
	}

	c.JSON(http.StatusOK, &utils.ApiMessageHandler{
		Success: true,
		StatusCode: http.StatusOK,
		Message: "funds transferred to wallet successfully",
		Data: map[string]any{"amount": payload.Amount},
		Error: nil,
	})
}