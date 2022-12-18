package database

import (
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var RedisInstance *redis.Client
