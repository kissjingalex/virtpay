package config

import (
	"github.com/kissjingalex/virtpay/internal/util/db/postgresdb"
	gredisconf "github.com/kissjingalex/virtpay/internal/util/db/redis"
)

var (
	// GlobalConfig config
	GlobalConfig = &Config{}
)

type Config struct {
	ServiceName string
	Env         string
	MDB         *postgresdb.PostgresConfig
	MRedis      *gredisconf.RedisConfig
}
