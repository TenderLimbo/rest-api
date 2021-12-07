package repository

import (
	"github.com/TenderLimbo/rest-api/models"
	"gorm.io/gorm"
)

type BooksManager interface {
	GetBooks(filterCondition map[string][]string) ([]models.Book, error)
	GetBookByID(id int) (models.Book, error)
	CreateBook(book models.Book) (int, error)
	DeleteBookByID(id int) error
	UpdateBookByID(id int, book models.Book) error
}

type BooksManagerPostgres struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *BooksManagerPostgres {
	return &BooksManagerPostgres{db: db}
}

func (r *BooksManagerPostgres) GetBooks(filterCondition map[string][]string) ([]models.Book, error) {
	var Books []models.Book
	var err error
	if len(filterCondition) == 0 {
		err = r.db.Where("amount <> ?", 0).Find(&Books).Error
	} else {
		err = r.db.Where("amount <> ?", 0).Where("genre = ?", filterCondition["genre"]).Find(&Books).Error
	}
	return Books, err
}

func (r *BooksManagerPostgres) GetBookByID(id int) (models.Book, error) {
	var book models.Book
	if err := r.db.First(&book, id).Error; err != nil {
		return book, err
	}
	return book, nil
}

func (r *BooksManagerPostgres) CreateBook(newBook models.Book) (int, error) {
	if err := r.db.Debug().Select("name", "price", "genre", "amount").Create(&newBook).Error; err != nil {
		return newBook.ID, err
	}
	return newBook.ID, nil
}

func (r *BooksManagerPostgres) DeleteBookByID(id int) error {
	res := r.db.Delete(&models.Book{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BooksManagerPostgres) UpdateBookByID(id int, newBook models.Book) error {
	res := r.db.Where("id = ?", id).Select("*").Omit("id").Updates(newBook)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
