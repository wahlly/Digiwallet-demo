package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusOK, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: map[string]any{},
			Error: err,
		})
		return
	}

	id, token, err := uc.UserService.LoginUser(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ApiMessageHandler{
			Success: false,
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: map[string]any{},
			Error: err,
		})
		return
	}

	c.JSON(http.StatusOK, utils.ApiMessageHandler{
		Success: true,
		StatusCode: http.StatusOK,
		Message: "success",
		Data: map[string]any{"userId": id, "token": token},
		Error: nil,
	})
}