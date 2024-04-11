package pg

import (
	"github.com/kissjingalex/virtpay/internal/util/db/postgresdb"
	"github.com/kissjingalex/virtpay/internal/wallet/config"
	"gorm.io/gorm"
)

type WalletRepos struct {
	mDB *gorm.DB
}

func NewWalletRepos() *WalletRepos {
	mdb, err := postgresdb.DB(config.GlobalConfig.MDB)
	if err != nil {
		panic(err)
	}

	return &WalletRepos{
		mDB: mdb,
	}
}
