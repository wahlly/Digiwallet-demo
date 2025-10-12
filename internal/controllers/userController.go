package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/internal/dtos"
	"github.com/wahlly/Digiwallet-demo/internal/models"
	"github.com/wahlly/Digiwallet-demo/internal/services"
	"github.com/wahlly/Digiwallet-demo/internal/utils"
)

type UserController struct{
	UserService services.UserService
}

func (uc *UserController) CreateUser(c *gin.Context) {
	reqBody, err := utils.BindAndValidateReqBody[dtos.UserRegistrationReqDto](c)
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

	user := &models.User{
		FirstName: reqBody.FirstName,
		LastName: reqBody.LastName,
		Phone: reqBody.Phone,
		Email: reqBody.Email,
		Password: reqBody.Password,
		UserName: reqBody.UserName,
	}
	err = uc.UserService.CreateUser(ctx, user)
	if err != nil {
		fmt.Println(err)
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
		Message: "User created successfully",
		Data: map[string]any{},
	})
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