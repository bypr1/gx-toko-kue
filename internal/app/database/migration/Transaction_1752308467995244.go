package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type Transaction_1752308467995244 struct{}

func (Transaction_1752308467995244) Reference() string {
	return "Transaction_1752308467995244"
}

func (Transaction_1752308467995244) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")

	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.Transaction{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: model.TransactionCake{}, Owner: owner},
	}
}

func (Transaction_1752308467995244) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
