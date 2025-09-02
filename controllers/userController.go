package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/dtos"
	"github.com/wahlly/Digiwallet-demo/models"
	"github.com/wahlly/Digiwallet-demo/services"
	"github.com/wahlly/Digiwallet-demo/utils"
)

type UserController struct{
	UserService services.UserService
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := uc.UserService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) LoginUser(c *gin.Context) {
	reqBody, err := utils.BindAndValidateReqBody[dtos.UserLoginReqDto](c)
	if err != nil {
		e := utils.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusBadRequest,
			Message: "validation error",
			Data: map[string]any{
				"errors": e,
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, token, err := uc.UserService.LoginUser(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: map[string]any{},
		})
		return
	}

	c.JSON(http.StatusOK, utils.ApiMessageHandler{
		Success: true,
		StatusCode: http.StatusOK,
		Message: "success",
		Data: map[string]any{"userId": id, "token": token},
	})
}