package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/chihabMe/jwt-refresh-token/models"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateTokens(user models.User) (map[string]string, error) {

	accessTime, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TIME"))
	refreshTime, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIME"))

	token := jwt.New(jwt.SigningMethodHS256)
	accessClaims := token.Claims.(jwt.MapClaims)
	accessClaims["user_id"] = user.ID
	accessClaims["username"] = user.Username
	accessClaims["exp"] = time.Now().Add(time.Hour * time.Duration(accessTime)).Unix()
	//refresh token
	refresh := jwt.New(jwt.SigningMethodES256)
	refreshClaims := refresh.Claims.(jwt.MapClaims)
	refreshClaims["user_id"] = user.ID
	refreshClaims["exp"] = time.Now().Add(time.Hour * time.Duration(refreshTime)).Unix()
	//acc, err := token.SignedString([]byte(os.Getenv(("SECRET"))))
	acc, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, err
	}
	ref, err := token.SignedString([]byte(os.Getenv(("SECRET"))))
	if err != nil {
		return nil, err
	}
	return map[string]string{"access_token": acc, "refresh_token": ref}, nil
}

func VerifyTokenValidity(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method")
		}
		return []byte(os.Getenv("SECRET")), nil

	})
	if err != nil {
		return nil, err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	alive := VerifyTokenExpireDate(token)
	if !alive {
		return nil, errors.New("dead token")
	}
	return token, nil

}
func VerifyTokenExpireDate(token *jwt.Token) bool {
	claims := token.Claims.(jwt.MapClaims)
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return false
	}
	return true

}
