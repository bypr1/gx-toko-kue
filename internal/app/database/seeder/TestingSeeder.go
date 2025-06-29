package seeder

import (
	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type TestingSeeder struct{}

func (seed *TestingSeeder) Seed() {
	testings := seed.setTestingData()
	for _, testing := range testings {
		var count int64
		config.PgSQL.Model(&model.Testing{}).Where("name = ?", testing["name"]).Count(&count)
		if count > 0 {
			continue
		}

		config.PgSQL.Create(&model.Testing{
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
