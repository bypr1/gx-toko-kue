package migration

import (
	"service/internal/pkg/config"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type CakeTableRename_1752237331222387 struct{}

func (CakeTableRename_1752237331222387) Reference() string {
	return "CakeTableRename_1752237331222387"
}

func (CakeTableRename_1752237331222387) Tables() []xtremedb.Table {
	return []xtremedb.Table{
		{
			Connection: config.PgSQL,
			RenameTable: xtremedb.Rename{
				Old: "cake_recipes",
				New: "cake_recipe_ingredients",
			},
		},
		{
			Connection: config.PgSQL,
			RenameTable: xtremedb.Rename{
				Old: "cake_ingredients",
				New: "cake_component_ingredients",
			},
		},
		{
			Connection: config.PgSQL,
			RenameTable: xtremedb.Rename{
				Old: "cake_cakes",
				New: "cakes",
			},
		},
	}
}

func (CakeTableRename_1752237331222387) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
