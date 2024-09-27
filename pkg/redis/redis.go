package redis

import (
	"github.com/gofiber/storage/redis"
	"micro-fiber-test/pkg/config"
	"runtime"
)

func ConfigureRedisStorage(configuration *config.Configuration) *redis.Storage {
	return redis.New(redis.Config{
		Host:      configuration.RedisHost,
		Port:      configuration.RedisPort,
		Username:  configuration.RedisUser,
		Password:  configuration.RedisPass,
		URL:       "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	},
	)
}
