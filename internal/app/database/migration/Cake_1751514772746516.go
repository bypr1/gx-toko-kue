package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type Cake_1751514772746516 struct{}

func (Cake_1751514772746516) Reference() string {
	return "Cake_1751514772746516"
}

func (s Cake_1751514772746516) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")

	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.CakeComponentIngredient{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: model.Cake{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: model.CakeIngredient{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: model.CakeCost{}, Owner: owner},
	}
}

func (Cake_1751514772746516) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
