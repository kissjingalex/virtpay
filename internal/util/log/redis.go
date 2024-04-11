package log

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"time"

	"code.cloudfoundry.org/go-diodes"
)

// RedisConfig is conf for redis writer.
type RedisConfig struct {
	Level     Level
	RedisURL  string
	RedisPass string
	LogKey    string
	client    *redis.Client
}

// RedisWriter generate redis writer.
func RedisWriter(redisConf ...RedisConfig) io.Writer {
	conf := redisConf[0]
	if len(redisConf) == 0 {
		conf = RedisConfig{
			Level: InfoLevel,
		}
	}
	if conf.RedisURL == "" || conf.RedisPass == "" {
		conf.RedisURL = "log-proxy.ops:6379"
		conf.RedisPass = ""
	}
	if conf.LogKey == "" {
		conf.LogKey = "wtserver:basic:log"
	}
	conf.client = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     conf.RedisURL,
		Password: conf.RedisPass,
	})
	if async {
		w := NewAsyncWriter(conf.Level, conf, diodes.NewManyToOne(1024, diodes.AlertFunc(func(missed int) {
			fmt.Printf("Redis dropped %d messages\n", missed)
		})), 1*time.Second)
		registerCloseFunc(w.Close)
		return w
	}
	return conf
}

// Write write data to writer
func (c RedisConfig) Write(p []byte) (n int, err error) {
	return len(p), c.client.LPush(context.Background(), c.LogKey, p).Err()
}

// WriteLevel write data to writer with level info provided
func (c RedisConfig) WriteLevel(level Level, p []byte) (n int, err error) {
	if level < c.Level {
		return len(p), nil
	}

	return c.Write(p)
}
