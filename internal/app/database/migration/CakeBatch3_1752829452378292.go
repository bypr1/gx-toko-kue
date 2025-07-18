package migration

import (
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type CakeBatch3_1752829452378292 struct{}

func (CakeBatch3_1752829452378292) Reference() string {
	return "CakeBatch3_1752829452378292"
}

func (CakeBatch3_1752829452378292) Tables() []xtremedb.Table {
	return []xtremedb.Table{}
}

func (CakeBatch3_1752829452378292) Columns() []xtremedb.Column {
	return []xtremedb.Column{
		{
			Connection: config.PgSQL,
			Model:      model.Cake{},
			RenameColumns: []xtremedb.Rename{
				{
					Old: "sellPrice",
					New: "price",
				},
			},
		},
	}
}
