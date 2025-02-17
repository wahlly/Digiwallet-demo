package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/config"
	"github.com/wahlly/Digiwallet-demo/routes"
)

func main() {
	router := gin.New()
	db := config.ConnectDB()

	routes.RegisterRoutes(router, db)
	router.GET("/", func(ctx *gin.Context) {ctx.String(http.StatusOK, "Hello Gopher!")})

	router.Run(":6040")
}