// @Title
// @Description
// @Author
// @Update
package handlers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/chihabMe/jwt-refresh-token/database"
	"github.com/chihabMe/jwt-refresh-token/models"
	"github.com/chihabMe/jwt-refresh-token/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func ObtainToken(c *fiber.Ctx) error {
	type Inputs struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var inputs Inputs
	if err := c.BodyParser(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "please check your credentials"})
	}
	var user models.User = models.User{Email: inputs.Email}
	if err := database.Instance.Find(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "please check your credentials"})
	}
	same := utils.CheckPassword(inputs.Password, user)
	if !same {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "please check your credentials"})
	}
	tokens, err := utils.GenerateTokens(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "please check your credentials"})
	}
	var ctx context.Context = context.Background()
	utils.SetRefreshToken(ctx, user.ID, tokens["refresh_token"])
	ctx.Done()
	accessTime, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TIME"))
	refreshTime, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIME"))
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokens["access_token"],
		Path:     "/",
		MaxAge:   int(time.Now().Add(time.Duration(accessTime) * time.Hour).Unix()),
		Secure:   false,
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    tokens["refresh_token"],
		Path:     "/",
		MaxAge:   int(time.Now().Add(time.Duration(refreshTime) * time.Hour).Unix()),
		Secure:   false,
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": tokens})
}
func RefreshToken(c *fiber.Ctx) error {
	refreshString := c.Cookies("refresh")
	token, err := utils.VerifyTokenValidity(refreshString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "dead refresh token please login again"})
	}
	claims := token.Claims.(jwt.MapClaims)
	//getting the user from the database
	// by extracting the user ID from the token
	var user models.User
	if err := database.Instance.Where("id=?", claims["user_id"]).Find(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "dead refresh token please login again"})
	}
	//checking if the received refresh token is the latest generated refresh token for that user
	//by  looking in the redis database
	ctx := context.Background()
	latestRefreshToken, err := utils.GetRefreshToken(ctx, user.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "dead refresh token please login again"})
	}
	sameRefreshToken := latestRefreshToken == refreshString
	if !sameRefreshToken {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "dead refresh token please login again"})
	}
	// generate a new refresh/access token
	// store the new refresh token in redis
	// set the new cookies
	tokens, err := utils.GenerateTokens(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "dead refresh token please login again"})
	}
	//utils.DeleteRefreshToken(ctx,user.ID)
	err = utils.SetRefreshToken(ctx, user.ID, tokens["refresh_token"])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "dead refresh token please login again"})
	}
	accessTime, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TIME"))
	refreshTime, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIME"))
	ctx.Done()
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokens["access_token"],
		Path:     "/",
		MaxAge:   int(time.Now().Add(time.Duration(accessTime) * time.Hour).Unix()),
		Secure:   false,
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    tokens["refresh_token"],
		Path:     "/",
		MaxAge:   int(time.Now().Add(time.Duration(refreshTime) * time.Hour).Unix()),
		Secure:   false,
		HTTPOnly: true,
	})

	return c.SendStatus(fiber.StatusOK)
}
func LogoutToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh")
	fmt.Println("pass 1")
	fmt.Println("refresh stoken:", refreshToken)
	token, err := utils.VerifyTokenValidity(refreshToken)
	if err != nil {
		//fake ok status
		return c.SendStatus(fiber.StatusOK)
		//return c.SendStatus(fiber.StatusBadRequest)
	}
	fmt.Println("pass 2")
	claims := token.Claims.(jwt.MapClaims)
	var user models.User
	userId := claims["user_id"]
	if err := database.Instance.Where("id=?", userId).Find(&user).Error; err != nil {
		//fake ok status
		return c.SendStatus(fiber.StatusOK)
		//return c.SendStatus(fiber.StatusBadRequest)
	}
	ctx := context.Background()
	latestRefreshToken, err := utils.GetRefreshToken(ctx, user.ID)
	if err != nil {
		//fake ok status
		return c.SendStatus(fiber.StatusOK)
		//return c.SendStatus(fiber.StatusBadRequest)
	}
	sameRefreshToken := latestRefreshToken == refreshToken
	if !sameRefreshToken {
		//fake ok status
		return c.SendStatus(fiber.StatusOK)
		//return c.SendStatus(fiber.StatusBadRequest)
	}
	utils.DeleteRefreshToken(ctx, user.ID)
	return c.SendStatus(fiber.StatusOK)
}
func VerifyToken(c *fiber.Ctx) error {
	tokenString := c.Cookies("Authorization")
	_, err := utils.VerifyTokenValidity(tokenString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "invalid token "})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: true,
	})
	return c.SendStatus(fiber.StatusOK)
}
