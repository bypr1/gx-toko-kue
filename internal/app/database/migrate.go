package database

import (
	"service/internal/app/database/migration"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"gorm.io/gorm"
)

func Migrations(conn *gorm.DB) []xtremedb.Migration {
	return []xtremedb.Migration{
		&migration.Activity_1726651211960757{Connection: conn},
		&migration.Cake_1751514772746516{Connection: conn},
		&migration.CakeTableRename_1752237331222387{},
		&migration.Transaction_1752308467995244{},
		&migration.Cake3AddImage_1752565306321794{},
	}
}
