package database

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		panic(err)
	}
	Instance = db
	Migrate()
}
func ConnectRedis() {
	redisUrl := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})
	RedisInstance = rdb
	ctx := context.Background()
	err := rdb.Set(ctx, "key", "chihab eddin", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key:", val)

}
