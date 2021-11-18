package REST_API

import "errors"

type Model struct {
	ID int `json:"id"`
}

type Book struct {
	Model
	Name   string  `json:"name" gorm:"unique;not null"`
	Price  float64 `json:"price" gorm:"not null"`
	Genre  int     `json:"genre" gorm:"not null"`
	Amount int     `json:"amount" gorm:"not null"`
}

func ValidateBook(book Book) error {
	if len(book.Name) > 0 && len(book.Name) < 100 && book.Price >= 0 && book.Genre > 0 && book.Genre < 4 && book.Amount >= 0 {
		return nil
	}
	return errors.New("book is invalid")
}

type Genre struct {
	Model
	Name string
}
