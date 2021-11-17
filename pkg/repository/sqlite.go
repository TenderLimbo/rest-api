package repository

import (
	restapi "github.com/TenderLimbo/rest-api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqliteDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("Books.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&restapi.Book{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
