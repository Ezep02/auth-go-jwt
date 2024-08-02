package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("error hashing the password %v", err)
	}

	return string(hash), nil
}

func CheckPassword(password string, hashedPassword string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
