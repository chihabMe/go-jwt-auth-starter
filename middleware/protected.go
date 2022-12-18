package middleware

import (
	"github.com/chihabMe/jwt-refresh-token/database"
	"github.com/chihabMe/jwt-refresh-token/models"
	"github.com/chihabMe/jwt-refresh-token/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("Authorization")
		token, err := utils.VerifyTokenValidity(tokenString)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		claims := token.Claims.(jwt.MapClaims)
		var user models.User
		if err := database.Instance.Where("id=?", claims["user_id"]).Find(&user).Error; err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		c.Locals("user", user)
		return c.Next()
	}
}
