package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitRedis(cfg *viper.Viper, log *logrus.Logger) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.GetString("redis.host"), cfg.GetString("redis.port"))
	password := cfg.GetString("redis.password")
	db := cfg.GetInt("redis.db")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to Redis")
		return nil
	}

	log.Info("Connected to Redis successfully")
	return rdb
}
