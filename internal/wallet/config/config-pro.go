//go:build pro
// +build pro

package config

import (
	"github.com/kissjingalex/virtpay/internal/util/db"
	"github.com/kissjingalex/virtpay/internal/util/db/postgresdb"
	gredisconf "github.com/kissjingalex/virtpay/internal/util/db/redis"
)

func init() {
	GlobalConfig = &Config{
		ServiceName: "service-wallet",
		Env:         "pro",
		MDB: &postgresdb.PostgresConfig{
			DSN:           "",
			MaxOpenConns:  128,
			MaxIdleConns:  32,
			MaxLifetime:   300,
			MaxIdleTime:   120,
			LogLevel:      db.LogLevelWarn,
			SlowThreshold: 1,
		},
		MRedis: &gredisconf.RedisConfig{
			Addr:               "10.41.158.197:6379",
			Password:           "71f13506-6f7d-459b-a663-a38f9bc51891",
			MinIdleConns:       32,
			PoolSize:           64,
			MaxRetries:         1,
			ConnMaxIdleSecond:  300,
			ConnMaxLifeSecond:  300,
			ReadTimeoutSecond:  1,
			WriteTimeoutSecond: 1,
			DB:                 8,
			TLS:                false,
		},
	}

	GlobalConfig.MDB.LoadDSN("db_wallet")
}
