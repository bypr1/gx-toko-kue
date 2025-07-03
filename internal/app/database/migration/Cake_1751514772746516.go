package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model/cake"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type Cake_1751514772746516 struct{}

func (Cake_1751514772746516) Reference() string {
	return "Cake_1751514772746516"
}

func (Cake_1751514772746516) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")

	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: cake.Ingredient{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: cake.Cake{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: cake.CakeRecipe{}, Owner: owner},
		{Connection: config.PgSQL, CreateTable: cake.CakeCost{}, Owner: owner},
	}
}

func (Cake_1751514772746516) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
