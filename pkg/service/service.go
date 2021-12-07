package service

import (
	"github.com/TenderLimbo/rest-api/models"
	"github.com/TenderLimbo/rest-api/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type BooksManager interface {
	GetBooks(filterCondition map[string][]string) ([]models.Book, error)
	GetBookByID(id int) (models.Book, error)
	CreateBook(book models.Book) (int, error)
	DeleteBookByID(id int) error
	UpdateBookByID(id int, book models.Book) error
}

type BooksManagerService struct {
	repo repository.BooksManager
}

func NewService(repo repository.BooksManager) *BooksManagerService {
	return &BooksManagerService{repo: repo}
}

func (s *BooksManagerService) CreateBook(book models.Book) (int, error) {
	return s.repo.CreateBook(book)
}

func (s *BooksManagerService) GetBookByID(id int) (models.Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *BooksManagerService) GetBooks(filterCondition map[string][]string) ([]models.Book, error) {
	return s.repo.GetBooks(filterCondition)
}

func (s *BooksManagerService) DeleteBookByID(id int) error {
	return s.repo.DeleteBookByID(id)
}

func (s *BooksManagerService) UpdateBookByID(id int, book models.Book) error {
	return s.repo.UpdateBookByID(id, book)
}
