package utils

import (
	"github.com/chihabMe/jwt-refresh-token/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 15)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(pass string, user models.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	return err == nil
}
