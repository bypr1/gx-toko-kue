package database

import (
	"service/internal/app/database/migration"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"gorm.io/gorm"
)

func Migrations(conn *gorm.DB) []xtremedb.Migration {
	return []xtremedb.Migration{
		&migration.Activity_1726651211960757{},
		&migration.Cake_1751514772746516{},
		&migration.Transaction_1752308467995244{},
		&migration.CakeBatch2_1752829452378290{},
		&migration.TransactionBatch2_1752829452378291{},
		&migration.CakeBatch3_1752829452378292{},
	}
}
