package config

import (
	"os"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"gorm.io/gorm"
)

var (
	PgSQL *gorm.DB
)

func InitDB() func() {
	var DBClose func()

	conf := xtremedb.DBConf{
		Driver:    xtremedb.POSTGRESQL_DRIVER,
		Host:      os.Getenv("DB_HOST"),
		Port:      os.Getenv("DB_PORT"),
		Username:  os.Getenv("DB_USERNAME"),
		Password:  os.Getenv("DB_PASSWORD"),
		Database:  os.Getenv("DB_DATABASE"),
		ParseTime: true,
	}

	PgSQL, DBClose = xtremedb.Connect(conf)

	return DBClose
}
