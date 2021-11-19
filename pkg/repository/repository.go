package repository

import (
	"errors"
	restapi "github.com/TenderLimbo/rest-api"
	"gorm.io/gorm"
)

type BooksManager interface {
	GetBooks(filterCondition map[string][]string) ([]restapi.Book, error)
	GetBookByID(id int) (restapi.Book, error)
	CreateBook(book restapi.Book) (int, error)
	DeleteBookByID(id int) error
	UpdateBookByID(id int, book restapi.Book) error
}

type BooksManagerPostgres struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *BooksManagerPostgres {
	return &BooksManagerPostgres{db: db}
}

func (r *BooksManagerPostgres) GetBooks(filterCondition map[string][]string) ([]restapi.Book, error) {
	var Books []restapi.Book
	if len(filterCondition) == 0 {
		r.db.Where("amount <> ?", 0).Find(&Books)
	} else {
		r.db.Where("amount <> ?", 0).Where("genre = ?", filterCondition["genre"]).Find(&Books)
	}
	return Books, nil
}

func (r *BooksManagerPostgres) GetBookByID(id int) (restapi.Book, error) {
	var book restapi.Book
	if err := r.db.First(&book, id).Error; err != nil {
		return book, err
	}
	return book, nil
}

func (r *BooksManagerPostgres) CreateBook(newBook restapi.Book) (int, error) {
	if err := r.db.Create(&newBook).Error; err != nil {
		return newBook.ID, err
	}
	return newBook.ID, nil
}

func (r *BooksManagerPostgres) DeleteBookByID(id int) error {
	if r.db.Delete(&restapi.Book{}, id).RowsAffected < 1 {
		return errors.New("id not found")
	}
	return nil
}

func (r *BooksManagerPostgres) UpdateBookByID(id int, newBook restapi.Book) error {
	res := r.db.Where("id = ?", id).Select("*").Omit("id").Updates(newBook)
	if res.Error != nil || res.RowsAffected < 1 {
		return errors.New("id not found")
	}
	return nil
}
