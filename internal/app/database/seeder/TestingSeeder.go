package seeder

import (
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

type TestingSeeder struct {
	Connection *gorm.DB
}

func (seed *TestingSeeder) Seed() {
	testings := seed.setTestingData()
	for _, testing := range testings {
		var count int64
		seed.Connection.Model(&model.Testing{}).Where("name = ?", testing["name"]).Count(&count)
		if count > 0 {
			continue
		}

		seed.Connection.Create(&model.Testing{
			Name: testing["name"].(string),
		})
	}
}

/** --- UNEXPORTED FUNCTIONS --- */

func (seed *TestingSeeder) setTestingData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name": "Testing",
		},
	}
}
