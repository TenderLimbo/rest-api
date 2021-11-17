package service

import (
	restapi "github.com/TenderLimbo/rest-api"
	"github.com/TenderLimbo/rest-api/pkg/repository"
)

type BooksManager interface {
	GetBooks(filterCondition map[string][]string) ([]restapi.Book, error)
	GetBookByID(id string) (restapi.Book, error)
	CreateBook(book restapi.Book) (int, error)
	DeleteBookByID(id string) error
	UpdateBookByID(id string, book restapi.Book) error
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

func (s *BooksManagerService) GetBookByID(id string) (restapi.Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *BooksManagerService) GetBooks(filterCondition map[string][]string) ([]restapi.Book, error) {
	return s.repo.GetBooks(filterCondition)
}

func (s *BooksManagerService) DeleteBookByID(id string) error {
	return s.repo.DeleteBookByID(id)
}

func (s *BooksManagerService) UpdateBookByID(id string, book restapi.Book) error {
	return s.repo.UpdateBookByID(id, book)
}
