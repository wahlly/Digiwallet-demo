package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wahlly/Digiwallet-demo/models"
)


type SignedDetails struct{
	Email	string
	id	uint
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateJWTtoken(user models.User) (signedToken string, err error) {
	claims := &SignedDetails{
		Email: user.Email,
		id: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateJWTtoken(token string) (valid bool, claims *SignedDetails, msg string) {
	tk, err := jwt.ParseWithClaims(token, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return false, nil, err.Error()
	}

	claims, ok := tk.Claims.(*SignedDetails)
	if !ok {
		return false, nil, "token is invalid"
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return false, nil, "token is already expired"
	}

	return true, claims, "token is valid"
}