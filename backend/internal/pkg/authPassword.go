package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

const salt = bcrypt.DefaultCost

func HashPassword(password string) (string, error) {
	convertPass := []byte(password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(convertPass, salt)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ValidatePassword(password string, hashedPassword []byte) (bool, error) {
	convertPass := []byte(password)

	// To compare the password with the hash
	err := bcrypt.CompareHashAndPassword(hashedPassword, convertPass)

	if err != nil {
		return false, err
	}
	return true, nil
}
