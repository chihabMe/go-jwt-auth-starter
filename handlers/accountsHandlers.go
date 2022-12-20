package handlers

import (
	"github.com/chihabMe/jwt-refresh-token/database"
	"github.com/chihabMe/jwt-refresh-token/models"
	"github.com/chihabMe/jwt-refresh-token/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	return c.JSON("users list")
}
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := database.Instance.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": err})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": users})
}

func DeleteUser(c *fiver.Ctx){
	return c.JSON(fiber.Map{"status":"success","data":"deleted"})
}

func Register(c *fiber.Ctx) error {
	type Response struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": err})
	}
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "can't register an new account right now"})
	}
	user.Password = hash
	if err := database.Instance.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "data": "please make sure that you are using a unique username and email"})
	}
	response := Response{Username: user.Username, Email: user.Email}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": response})
}
