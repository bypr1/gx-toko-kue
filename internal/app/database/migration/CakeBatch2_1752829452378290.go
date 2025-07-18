package migration

import (
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type CakeBatch2_1752829452378290 struct{}

func (CakeBatch2_1752829452378290) Reference() string {
	return "CakeBatch2_1752829452378290"
}

func (CakeBatch2_1752829452378290) Tables() []xtremedb.Table {
	return []xtremedb.Table{}
}

func (CakeBatch2_1752829452378290) Columns() []xtremedb.Column {
	return []xtremedb.Column{
		{
			Connection: config.PgSQL,
			Model:      model.Cake{},
			RenameColumns: []xtremedb.Rename{
				{
					Old: "unit",
					New: "unitId",
				},
			},
			AlterColumns: []string{"description"},
		},
		{
			Connection: config.PgSQL,
			Model:      model.CakeComponentIngredient{},
			RenameColumns: []xtremedb.Rename{
				{
					Old: "unitPrice",
					New: "price",
				},
			},
			AlterColumns: []string{"description"},
		},
		{
			Connection: config.PgSQL,
			Model:      model.CakeCost{},
			RenameColumns: []xtremedb.Rename{
				{
					Old: "cost",
					New: "price",
				},
			},
		},
	}
}
