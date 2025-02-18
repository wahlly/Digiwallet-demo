package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/services"
)


type WalletController struct {
	WalletService *services.WalletService 
}

func (wc *WalletController) GetUserWallet(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, sign in again"})
		c.Abort()
		return
	}

	wallet, err := wc.WalletService.GetUserWallet(email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "wallet": wallet})
}