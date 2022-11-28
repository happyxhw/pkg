package goredis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
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
	Host         string
	Password     string
	DB           int
	PoolSize     int `mapstructure:"pool_size"`
	MinIdleConns int `mapstructure:"min_idle_conns"`
}

// NewRedis initialize redis instance
func NewRedis(redisConf *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisConf.Host,
		DB:           redisConf.DB,
		Password:     redisConf.Password,
		PoolSize:     redisConf.PoolSize,
		MinIdleConns: redisConf.MinIdleConns,
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
