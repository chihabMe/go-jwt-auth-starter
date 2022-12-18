// @Title
// @Description
// @Author
// @Update
package routes

import (
	"github.com/chihabMe/jwt-refresh-token/handlers"
	"github.com/chihabMe/jwt-refresh-token/middleware"
	fiber "github.com/gofiber/fiber/v2"
)

func RegisterAccounts(app *fiber.Router) {
	accountsRoutes := (*app).Group("accounts/")
	accountsRoutes.Post("register/", handlers.Register)
	accountsRoutes.Get("user/:id/", middleware.Protected(), handlers.GetUser)
	accountsRoutes.Get("users/", middleware.Protected(), handlers.GetAllUsers)
}
