package migration

import (
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type TransactionBatch2_1752829452378291 struct{}

func (TransactionBatch2_1752829452378291) Reference() string {
	return "TransactionBatch2_1752829452378291"
}

func (TransactionBatch2_1752829452378291) Tables() []xtremedb.Table {
	return []xtremedb.Table{}
}

func (TransactionBatch2_1752829452378291) Columns() []xtremedb.Column {
	return []xtremedb.Column{
		{
			Connection: config.PgSQL,
			Model:      model.Transaction{},
			RenameColumns: []xtremedb.Rename{
				{
					Old: "transactionDate",
					New: "date",
				},
				{
					Old: "totalAmount",
					New: "totalPrice",
				},
			},
		},
		{
			Connection: config.PgSQL,
			Model:      model.TransactionCake{},
			RenameColumns: []xtremedb.Rename{
				{
					Old: "unitPrice",
					New: "price",
				},
			},
		},
	}
}
