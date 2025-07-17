package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type Activity_1726651211960757 struct{}

func (Activity_1726651211960757) Reference() string {
	return "Activity_1726651211960757"
}

func (s Activity_1726651211960757) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")

	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.Activity{}, Owner: owner},
	}
}

func (Activity_1726651211960757) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
