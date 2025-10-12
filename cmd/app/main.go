package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/wahlly/Digiwallet-demo/config"
	"github.com/wahlly/Digiwallet-demo/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	router := gin.New()
	db := config.ConnectDB()

	routes.RegisterRoutes(router, db)
	router.GET("/", func(ctx *gin.Context) {ctx.String(http.StatusOK, "Hello Gopher!")})

	router.Run(":6040")
}