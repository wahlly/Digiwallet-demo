package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "no authorization header"})
			ctx.Abort()
			return
		}

		valid, claims, msg := ValidateJWTtoken(clientToken)
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": msg})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("id", claims.id)
		ctx.Next()
	}
}