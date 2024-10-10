package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateTokenAndSetCookie(userID any, w http.ResponseWriter) (string, error) {
	err := godotenv.Load("./internal/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mySigningKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		// SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &cookie)

	return tokenString, nil
}
