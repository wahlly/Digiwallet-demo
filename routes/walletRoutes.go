package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahlly/Digiwallet-demo/controllers"
	"github.com/wahlly/Digiwallet-demo/services"
	"github.com/wahlly/Digiwallet-demo/utils"
)


func WalletRoutes(rg *gin.RouterGroup, ws *services.WalletService) {
	router := rg.Group("/wallet")
	wc := &controllers.WalletController{WalletService: ws}

	router.GET("/", utils.Authenticate(), wc.GetUserWallet)
	router.GET("/:address", utils.Authenticate(), wc.GetWalletByAddress)
	router.POST("/deposit/initialize", utils.Authenticate(), wc.InitializeWalletDeposit)
	router.GET("/deposit/verify/:reference", utils.Authenticate(), wc.VerifyWalletDeposit)
	router.POST("/transfer", utils.Authenticate(), wc.TransferToWalletAddress)
}