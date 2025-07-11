package migration

import (
	"os"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
	"gorm.io/gorm"
)

type Cake_1751514772746516 struct {
	Connection *gorm.DB
}

func (Cake_1751514772746516) Reference() string {
	return "Cake_1751514772746516"
}

func (s Cake_1751514772746516) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")

	return []xtremedb.Table{
		{Connection: s.Connection, CreateTable: model.CakeComponentIngredient{}, Owner: owner},
		{Connection: s.Connection, CreateTable: model.Cake{}, Owner: owner},
		{Connection: s.Connection, CreateTable: model.CakeRecipeIngredient{}, Owner: owner},
		{Connection: s.Connection, CreateTable: model.CakeCost{}, Owner: owner},
	}
}

func (Cake_1751514772746516) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
