package repository

import (
	"fmt"
	restapi "github.com/TenderLimbo/rest-api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Genres = []restapi.Genre{
	{Name: "Adventure"},
	{Name: "Classics"},
	{Name: "Fantasy"},
}

func NewPostgresDB(config map[string]string, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config["host"], config["user"], password, config["dbname"], config["port"], config["sslmode"])
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
