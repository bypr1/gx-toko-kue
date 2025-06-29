package core

import (
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
	"service/internal/pkg/config"
	"time"
)

type TransactionRepository interface {
	SetTransaction(tx *gorm.DB)
}

type FirstIdRepository[M any] interface {
	FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) M
}

type FirstUUIDRepository[M any] interface {
	FirstByUUID(uuid string, args ...func(query *gorm.DB) *gorm.DB) M
}

type FindRepository[M any] interface {
	Find(parameter url.Values) []M
}

type PaginateRepository[M any] interface {
	Paginate(parameter url.Values) ([]M, interface{}, error)
}

type NumberPoolRepository interface {
	TransactionRepository
	TakenNumberPool(number ...string) string
}

// TODO: Re-enable this code after installing github.com/globalxtreme/go-identifier module (If you use GX Identifier for authorization)
//type EmployeeIdentifierRepository interface {
//	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)
//}

func GetIncrementMonthly(model interface{}) int64 {
	var totalData int64
	config.PgSQL.Unscoped().
		Where("EXTRACT(MONTH FROM \"createdAt\") = ?", time.Now().Month()).
		Model(&model).
		Count(&totalData)

	return totalData + 1
}

func Truncate(db *gorm.DB, tables ...schema.Tabler) {
	if len(tables) > 0 {
		for _, table := range tables {
			err := db.Exec(fmt.Sprintf("truncate table %s restart identity cascade", table.TableName()))
			if err != nil {
				xtremepkg.LogError(fmt.Sprintf("Truncate invalid: %v", err), false)
			}
		}
	}
}

func TakenNumberPool(repo NumberPoolRepository, tx *gorm.DB) (string, func()) {
	var numberPool string
	tx.Transaction(func(tx *gorm.DB) error {
		repo.SetTransaction(tx)
		numberPool = repo.TakenNumberPool()

		return nil
	})

	errFunc := func() {
		if r := recover(); r != nil {
			repo.TakenNumberPool(numberPool)
			panic(r)
		}
	}

	return numberPool, errFunc
}
