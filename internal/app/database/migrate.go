package database

import (
	"service/internal/app/database/migration"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

func Migrations() []xtremedb.Migration {
	return []xtremedb.Migration{
		&migration.Activity_1726651211960757{},
		&migration.Cake_1751514772746516{},
	}
}
