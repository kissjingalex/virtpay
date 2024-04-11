package gredisconf

import (
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	Addr               string
	Username           string
	Password           string
	MinIdleConns       int
	PoolSize           int
	MaxRetries         int
	ReadTimeoutSecond  int
	WriteTimeoutSecond int
	ConnMaxLifeSecond  int
	ConnMaxIdleSecond  int
	DB                 int
	TLS                bool
}

func (c *RedisConfig) ConvertToGRedisOpts() *redis.Options {
	ops := &redis.Options{
		Addr:            c.Addr,
		Username:        c.Username,
		Password:        c.Password,
		MinIdleConns:    c.MinIdleConns,
		PoolSize:        c.PoolSize,
		MaxRetries:      c.MaxRetries,
		ConnMaxIdleTime: time.Second * time.Duration(c.ConnMaxLifeSecond),
		ConnMaxLifetime: time.Second * time.Duration(c.ConnMaxIdleSecond),
		ReadTimeout:     time.Second * time.Duration(c.ReadTimeoutSecond),
		WriteTimeout:    time.Second * time.Duration(c.WriteTimeoutSecond),
		DB:              c.DB,
	}
	if c.TLS {
		ops.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return ops
}
