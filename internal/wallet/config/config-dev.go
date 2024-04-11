package config

import (
	"github.com/kissjingalex/virtpay/internal/util/db"
	"github.com/kissjingalex/virtpay/internal/util/db/postgresdb"
	gredisconf "github.com/kissjingalex/virtpay/internal/util/db/redis"
	"github.com/kissjingalex/virtpay/internal/util/log"
)

func init() {
	GlobalConfig = &Config{
		ServiceName: "service-wallet",
		Env:         "dev",
		MDB: &postgresdb.PostgresConfig{
			DSN:           "",
			MaxOpenConns:  8,
			MaxIdleConns:  8,
			MaxLifetime:   300,
			MaxIdleTime:   120,
			LogLevel:      db.LogLevelWarn,
			SlowThreshold: 1,
		},
		MRedis: &gredisconf.RedisConfig{
			Addr:               "127.0.0.1:6379",
			MinIdleConns:       16,
			PoolSize:           128,
			MaxRetries:         1,
			ConnMaxIdleSecond:  300,
			ConnMaxLifeSecond:  300,
			ReadTimeoutSecond:  1,
			WriteTimeoutSecond: 1,
			DB:                 8,
		},
	}

	GlobalConfig.MDB.LoadDSN("db_wallet")

	log.Info("config").Msgf("service=%s, nv=%s", GlobalConfig.ServiceName, GlobalConfig.Env)
}
