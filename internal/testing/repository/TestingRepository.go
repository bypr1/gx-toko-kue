package repository

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
)

/** --- INTERFACE --- */

type TestingRepository interface {
	core.TransactionRepository

	FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Testing
	Find(parameter url.Values) ([]model.Testing, interface{}, error)

	Store(form form.TestingForm) model.Testing
	Delete(testing model.Testing)

	AddSub(testing model.Testing, sub string) model.TestingSub
	DeleteSub(testingSub model.TestingSub)
}

func NewTestingRepository(args ...*gorm.DB) TestingRepository {
	repository := testingRepository{}
	if len(args) > 0 {
		repository.transaction = args[0]
	}

	return &repository
}

/** --- MAIN REPOSITORY --- */

type testingRepository struct {
	transaction *gorm.DB
}

func (repo *testingRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *testingRepository) FirstById(id any, args ...func(query *gorm.DB) *gorm.DB) model.Testing {
	var testing model.Testing

	query := config.PgSQL
	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&testing, "id = ?", id).Error
	if err != nil {
		error2.ErrXtremeTestingGet(err.Error())
	}

	return testing
}

func (repo *testingRepository) Find(parameter url.Values) ([]model.Testing, interface{}, error) {
	fromDate, toDate := core.SetDateRange(parameter)

	query := config.PgSQL.Preload("Subs").
		Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)

	if search := parameter.Get("search"); len(search) > 3 {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	testings, pagination, err := xtrememodel.Paginate(query.Order("id DESC"), parameter, model.Testing{})
	if err != nil {
		return nil, nil, err
	}

	return testings, pagination, nil
}

func (repo *testingRepository) Store(form form.TestingForm) model.Testing {
	testing := model.Testing{
		Name: form.Name,
	}

	err := repo.transaction.Create(&testing).Error
	if err != nil {
		error2.ErrXtremeTestingSave(err.Error())
	}

	return testing
}

func (repo *testingRepository) Delete(testing model.Testing) {
	err := repo.transaction.Delete(&testing).Error
	if err != nil {
		error2.ErrXtremeTestingDelete(err.Error())
	}
}

func (repo *testingRepository) AddSub(testing model.Testing, sub string) model.TestingSub {
	testingSub := model.TestingSub{
		TestingId: testing.ID,
		Name:      sub,
	}

	err := repo.transaction.Create(&testingSub).Error
	if err != nil {
		error2.ErrXtremeTestingSubSave(err.Error())
	}

	return testingSub
}

func (repo *testingRepository) DeleteSub(testingSub model.TestingSub) {
	err := repo.transaction.Delete(&testingSub).Error
	if err != nil {
		error2.ErrXtremeTestingSubDelete(err.Error())
	}
}
