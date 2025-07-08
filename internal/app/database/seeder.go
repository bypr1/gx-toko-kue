package database

import (
	"service/internal/app/database/seeder"

	"gorm.io/gorm"
)

type DatabaseSeeder interface {
	Seed()
}

type data struct {
	DatabaseSeeder
}

func Seeder(conn *gorm.DB) {
	seeders := []data{
		{&seeder.TestingSeeder{Connection: conn}},
		{&seeder.CakeIngredientSeeder{Connection: conn}},
	}

	for _, seed := range seeders {
		seed.Seed()
	}
}
