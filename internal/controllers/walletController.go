package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/internal/dtos"
	"github.com/wahlly/Digiwallet-demo/internal/services"
	"github.com/wahlly/Digiwallet-demo/internal/utils"
)


type WalletController struct {
	WalletService *services.WalletService
}

func (wc *WalletController) GetUserWallet(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: "unauthorized, sign in again",
			Data: map[string]any{},
		})
		c.Abort()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	wallet, err := wc.WalletService.GetUserWallet(ctx, id.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusInternalServerError,
			Message: err.Error(),
			Data: map[string]any{},
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, utils.ApiMessageHandler{
		Success: true,
		StatusCode: http.StatusOK,
		Message: "success",
		Data: map[string]any{"wallet": wallet},
	})
}

func (wc *WalletController) GetWalletByAddress(c *gin.Context) {
	address := c.Param("address")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	wallet, err := wc.WalletService.GetWalletByAddress(ctx, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusInternalServerError,
			Message: err.Error(),
			Data: map[string]any{},
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, utils.ApiMessageHandler{
		Success: true,
		StatusCode: http.StatusOK,
		Message: "success",
		Data: map[string]any{"wallet": wallet},
	})
}

func (wc *WalletController) InitializeWalletDeposit(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: "unauthorized, sign in again",
			Data: map[string]any{},
		})
		c.Abort()
		return
	}

	reqBody, validationErr := utils.BindAndValidateReqBody[dtos.InitTxnReqBody](c)
	if validationErr != nil {
		e := utils.FormatValidationErrors(validationErr)
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusBadRequest,
			Message: "validation error",
			Data: map[string]any{"errors": e},
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := wc.WalletService.InitializeWalletDeposit(ctx, reqBody, id.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: map[string]any{"errors": err},
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, utils.ApiMessageHandler{
		Success: true,
		StatusCode: http.StatusOK,
		Message: "success",
		Data: map[string]any{
			"uid": id,
			"payment_res": res,
		},
	})
}

func (wc *WalletController) VerifyWalletDeposit(c *gin.Context) {
	reference := c.Param("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusBadRequest,
			Message: "Invalid reference",
			Data: map[string]any{},
		})
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: "unauthorized, sign in again",
			Data: map[string]any{},
		})
		c.Abort()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx := wc.WalletService.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusInternalServerError,
			Message: "failed to start transaction",
			Data: map[string]any{},
		})
		c.Abort()
		return
	}
	res, err := wc.WalletService.VerifyWalletDeposit(ctx, tx, id.(uint), reference)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: err.Error(),
			Data: map[string]any{},
		})
		return
	}

	if err = tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: "failed to complete transaction",
			Data: map[string]any{},
		})
		return
	}

	c.JSON(http.StatusOK, utils.ApiMessageHandler{
		Success: false,
		StatusCode: http.StatusOK,
		Message: "wallet deposit successfull",
		Data: res,
	})
}

func (wc *WalletController) TransferToWalletAddress(c *gin.Context) {
	reqBody, err := utils.BindAndValidateReqBody[dtos.WalletP2PTransferReqDto](c)
	if err != nil {
		e := utils.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusBadRequest,
			Message: "validation error",
			Data: map[string]any{"errors": e},
		})
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: "unauthorized, sign in again",
			Data: map[string]any{},
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tx := wc.WalletService.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: "failed to start transaction",
			Data: map[string]any{},
		})
		return
	}
	defer func () {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, utils.ApiMessageHandler{
				Success: false,
				StatusCode: http.StatusUnauthorized,
				Message: "internal server error",
				Data: map[string]any{},
			})
			return
		}
	}()

	err = wc.WalletService.TransferToWalletAddress(ctx, tx, id.(uint), reqBody.Amount, reqBody.Recipient)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: err.Error(),
			Data: map[string]any{},
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusUnauthorized,
			Message: "failed to complete transaction",
			Data: map[string]any{},
		})
		return
	}

	c.JSON(http.StatusOK, &utils.ApiMessageHandler{
		Success: true,
		StatusCode: http.StatusOK,
		Message: "funds transferred to wallet successfully",
		Data: map[string]any{"amount": reqBody.Amount},
	})
}