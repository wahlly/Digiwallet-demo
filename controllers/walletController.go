package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/modules/paystack"
	"github.com/wahlly/Digiwallet-demo/services"
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