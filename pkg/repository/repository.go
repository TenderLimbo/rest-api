package repository

import (
	restapi "github.com/TenderLimbo/rest-api"
	"gorm.io/gorm"
)

type BooksManager interface {
	GetBooks(filterCondition map[string][]string) ([]restapi.Book, error)
	GetBookByID(id string) (restapi.Book, error)
	CreateBook(book restapi.Book) (int, error)
	DeleteBookByID(id string) error
	UpdateBookByID(id string, book restapi.Book) error
}

type BooksManagerSqlite struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *BooksManagerSqlite {
	return &BooksManagerSqlite{db: db}
}

func (r *BooksManagerSqlite) GetBooks(filterCondition map[string][]string) ([]restapi.Book, error) {
	var Books []restapi.Book
	if len(filterCondition) == 0 {
		r.db.Where("amount <> ?", 0).Find(&Books)
	} else if val, exists := filterCondition["genre"]; exists && len(val) == 1 {
		r.db.Where("amount <> ?", 0).Where("genre = ?", val).Find(&Books)
	} else {

	}
	return Books, nil
}

func (r *BooksManagerSqlite) GetBookByID(id string) (restapi.Book, error) {
	var book restapi.Book
	if err := r.db.First(&book, id).Error; err != nil {
		return book, err
	}
	return book, nil
}

func (r *BooksManagerSqlite) CreateBook(newBook restapi.Book) (int, error) {
	if err := r.db.Create(&newBook).Error; err != nil {
		return newBook.ID, err
	}
	return newBook.ID, nil
}

func (r *BooksManagerSqlite) DeleteBookByID(id string) error {
	if _, err := r.GetBookByID(id); err != nil {
		return err
	}
	r.db.Delete(&restapi.Book{}, id)
	return nil
}

func (r *BooksManagerSqlite) UpdateBookByID(id string, newBook restapi.Book) error {
	book, err := r.GetBookByID(id)
	if err != nil {
		return err
	}
	r.db.Model(&book).Select("*").Omit("id").Updates(newBook)
	return nil
}
