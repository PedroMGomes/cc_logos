package db

import (
	"go.search.crypto/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// pgsql:host=localhost;port=5432;dbname=testdb;user=bruce;password=mypass
// https://pkg.go.dev/gorm.io/driver/postgres?readme=expanded#section-readme
// Heroku provides the DB url via env variable DATABASE_URL.

// DB Database instance.
var DB *gorm.DB

// ConnectDatabase connects to a database and performs migrations for the models provided.
// If opening/connecting fails the process terminates with non-zero exit code (panic).
// If db migration fails the process terminates with non-zero exit code (panic).
func ConnectDatabase(databaseURL string) {
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN: databaseURL,
		// PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic("Failed to connect db: " + err.Error())
	}
	// note: AutoMigrate needs to be called for each model.
	err = database.AutoMigrate(&models.Currency{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
	DB = database
}
