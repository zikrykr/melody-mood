package setup

import (
	"context"

	"github.com/melody-mood/config"
	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func InitRedis() *redis.Client {
	conf := config.GetConfig()
	client := redis.NewClient(&redis.Options{
		Addr: conf.Redis.Host + ":" + conf.Redis.Port,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
	}
	logrus.Infof("Redis client initialized with address: %s:%s", conf.Redis.Host, conf.Redis.Port)
	return client
}
