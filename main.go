// @Title
// @Description
// @Author
// @Update
package main
//added text
import (
	"log"
	"os"

	"github.com/chihabMe/jwt-refresh-token/database"
	"github.com/chihabMe/jwt-refresh-token/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func RegisterRoutes(app *fiber.App) {
	app.Use(logger.New())
	v1 := app.Group("api/v1/")
	routes.RegisterAccounts(&v1)
	routes.RegisterAuth(&v1)
}
func main() {
	godotenv.Load()
	app := fiber.New()
	database.ConnectDB()
	database.ConnectRedis()
	RegisterRoutes(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
