package redis

import (
	"context"
	"errors"
	"fmt"
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

// NewClient new redis pool from config center by watch
func NewClient(redisConfig *Rdconfig) (*redis.Client, error) {
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

	err := redisCli.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("connect redis client failed %w", err)
	}

	return redisCli, nil
}
