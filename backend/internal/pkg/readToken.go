package pkg

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func ReadToken(cookieToken string) (string, error) {
	err := godotenv.Load("./internal/.env")
	if err != nil {
		return "", errors.New("error loading .env file")
	}

	mySigningKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(cookieToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if err != nil {
		return "", err
	}

	return token.Claims.(jwt.MapClaims)["userID"].(string), nil
}
