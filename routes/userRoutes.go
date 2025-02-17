package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/controllers"
	"github.com/wahlly/Digiwallet-demo/services"
)


func UserRoutes(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/user")
	uc := &controllers.UserController{UserService: userService}

	router.POST("/create", uc.CreateUser)
	router.POST("/login", uc.LoginUser)
}