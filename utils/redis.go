package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chihabMe/jwt-refresh-token/database"
	"github.com/google/uuid"
)

func SetRefreshToken(ctx context.Context, userId uuid.UUID, refreshToken string) error {
	refreshTime, _ := (strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIME")))

	//key := fmt.Sprintf("%s:%s", userId, refreshToken)
	key := fmt.Sprintf("%s", userId)
	if err := database.RedisInstance.Set(ctx, key, refreshToken, time.Duration(refreshTime)*time.Hour).Err(); err != nil {
		log.Println("Could not SET refresh token to redis for userID/tokenID")
		return errors.New("Could not SET refresh token to redis for userID/tokenID")
	}
	return nil
}

func DeleteRefreshToken(ctx context.Context, userId uuid.UUID) error {
	// key := fmt.Sprintf("%s:%s", userId, refreshToken)
	key := fmt.Sprintf("%s", userId)
	if err := database.RedisInstance.Del(ctx, key).Err(); err != nil {
		log.Println("Could not SET refresh token to redis for userID/tokenID")
		return errors.New("Could not delete refresh token to redis for userID/tokenID")
	}
	return nil
}
func GetRefreshToken(ctx context.Context, userId uuid.UUID) (string, error) {
	// key := fmt.Sprintf("%s:%s", userId, refreshToken)
	key := fmt.Sprintf("%s", userId)
	value, err := database.RedisInstance.Get(ctx, key).Result()
	if err != nil {
		log.Println("Could not SET refresh token to redis for userID/tokenID")
		return "", err
	}

	return value, nil
}
func GetRefreshTokenExpiration(ctx context.Context, userId uuid.UUID) time.Duration {
	// key := fmt.Sprintf("%s:%s", userId, refreshToken)
	key := fmt.Sprintf("%s", userId)
	value, err := database.RedisInstance.TTL(ctx, key).Result()
	if err != nil {
		log.Println("Could not SET refresh token to redis for userID/tokenID")
		return 0
	}

	return value
}
