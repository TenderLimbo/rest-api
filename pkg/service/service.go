package service

import (
	restapi "github.com/TenderLimbo/rest-api"
	"github.com/TenderLimbo/rest-api/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type BooksManager interface {
	GetBooks(filterCondition map[string][]string) ([]restapi.Book, error)
	GetBookByID(id int) (restapi.Book, error)
	CreateBook(book restapi.Book) (int, error)
	DeleteBookByID(id int) error
	UpdateBookByID(id int, book restapi.Book) error
}

type BooksManagerService struct {
	repo repository.BooksManager
}

func NewService(repo repository.BooksManager) *BooksManagerService {
	return &BooksManagerService{repo: repo}
}

func (s *BooksManagerService) CreateBook(book restapi.Book) (int, error) {
	return s.repo.CreateBook(book)
}

func (s *BooksManagerService) GetBookByID(id int) (restapi.Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *BooksManagerService) GetBooks(filterCondition map[string][]string) ([]restapi.Book, error) {
	return s.repo.GetBooks(filterCondition)
}

func (s *BooksManagerService) DeleteBookByID(id int) error {
	return s.repo.DeleteBookByID(id)
}

func (s *BooksManagerService) UpdateBookByID(id int, book restapi.Book) error {
	return s.repo.UpdateBookByID(id, book)
}
