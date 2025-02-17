package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/services"
	"gorm.io/gorm"
)

func RegisterRoutes(s *gin.Engine, db *gorm.DB) {
	basePath := s.Group("/v1")


	userService := services.NewUserService(db)
	UserRoutes(basePath, *userService)
}