// @Title
// @Description
// @Author
// @Update
package routes

import (
	"github.com/chihabMe/jwt-refresh-token/handlers"
	fiber "github.com/gofiber/fiber/v2"
)

func RegisterAuth(app *fiber.Router) {
	auth := (*app).Group("auth/")
	//obtain token || login
	auth.Post("token/obtain/", handlers.ObtainToken)
	//refresh token
	auth.Post("token/refresh/", handlers.RefreshToken)
	//verify
	auth.Post("token/verify/", handlers.VerifyToken)
	//logout
	auth.Post("token/logout/", handlers.LogoutToken)
}
