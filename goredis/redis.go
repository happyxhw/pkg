package goredis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/happyxhw/pkg/log"
)

var rdb *redis.Client

func InitDefaultRDB(cfg *Config) {
	var err error
	rdb, err = NewRedis(cfg)
	if err != nil {
		log.Fatal("init redis db", zap.Error(err))
	}
}

func DefaultRDB() *redis.Client {
	return rdb
}

// Config for go-redis
type Config struct {
	Addr         string
	Username     string
	Password     string
	DB           int
	PoolSize     int `mapstructure:"pool_size"`
	MinIdleConns int `mapstructure:"min_idle_conns"`
}

// NewRedis initialize redis instance
func NewRedis(redisConf *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		DB:           redisConf.DB,
		Username:     redisConf.Username,
		Password:     redisConf.Password,
		PoolSize:     redisConf.PoolSize,
		MinIdleConns: redisConf.MinIdleConns,
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
