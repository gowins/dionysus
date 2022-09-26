package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Rdconfig struct {
	MaxRetries  int    `json:"max_retries"`
	Addr        string `json:"addr"`
	DB          int    `json:"db"`
	Password    string `json:"password"`
	PoolSize    int    `json:"pool_size"`
	IdleTimeout int    `json:"idle_timeout"`
}

// new redis pool from config center by watch
func NewClients(redisConfig *Rdconfig) (*redis.Client, error) {
	if redisConfig.Addr == "" {
		return nil, errors.New("redis address is required")
	}

	if redisConfig.MaxRetries <= 0 {
		redisConfig.MaxRetries = 3
	}

	redisCli := redis.NewClient(&redis.Options{
		MaxRetries:  redisConfig.MaxRetries,
		Addr:        redisConfig.Addr,
		DB:          redisConfig.DB,
		PoolSize:    redisConfig.PoolSize,
		IdleTimeout: time.Duration(redisConfig.IdleTimeout) * time.Second,
		Password:    redisConfig.Password,
	})
	return redisCli, nil
}
