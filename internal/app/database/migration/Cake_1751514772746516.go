package migration

import (
	"os"
	"service/internal/pkg/model/cake"

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
		{Connection: s.Connection, CreateTable: cake.Ingredient{}, Owner: owner},
		{Connection: s.Connection, CreateTable: cake.CakeRecipe{}, Owner: owner},
		{Connection: s.Connection, CreateTable: cake.CakeCost{}, Owner: owner},
		{Connection: s.Connection, CreateTable: cake.Cake{}, Owner: owner},
	}
}

func (Cake_1751514772746516) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
