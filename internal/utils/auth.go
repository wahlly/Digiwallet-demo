package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)


func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerAuth := ctx.Request.Header.Get("authorization")
		auth := strings.Split(bearerAuth, " ")
		if len(auth) < 2 || auth[1] == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "no authorization header"})
			ctx.Abort()
			return
		}
		clientToken := auth[1]

		valid, claims, msg := ValidateJWTtoken(clientToken)
		if !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": msg})
			ctx.Abort()
			return
		}
		
		ctx.Set("email", claims.Email)
		ctx.Set("id", uint(claims.Id))
		ctx.Next()
	}
}