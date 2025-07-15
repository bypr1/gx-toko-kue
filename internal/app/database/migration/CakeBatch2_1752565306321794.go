package migration

import (
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type Cake3AddImage_1752565306321794 struct{}

func (Cake3AddImage_1752565306321794) Reference() string {
	return "CakeBatch2_1752565306321794"
}

func (Cake3AddImage_1752565306321794) Tables() []xtremedb.Table {
	return []xtremedb.Table{}
}

func (Cake3AddImage_1752565306321794) Columns() []xtremedb.Column {
	return []xtremedb.Column{
		{
			Connection: config.PgSQL,
			Model:      model.Cake{},
			AddColumns: []string{"image"},
		},
	}
}
