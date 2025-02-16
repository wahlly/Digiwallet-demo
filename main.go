package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/config"
)

func main() {
	router := gin.New()
	config.ConnectDB()

	router.Run(":6040")
}