package database

import (
	"log"

	"github.com/chihabMe/jwt-refresh-token/models"
)

func Migrate() {
	Instance.AutoMigrate(models.User{})
	log.Println("user model migrations=>done")
}
