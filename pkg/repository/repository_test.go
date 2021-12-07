package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TenderLimbo/rest-api/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func MockDB() (BooksManagerPostgres, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return BooksManagerPostgres{db: nil}, nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return BooksManagerPostgres{db: nil}, nil, err
	}
	return BooksManagerPostgres{db: gormDB}, mock, nil
}

func TestCreateBook(t *testing.T) {
	repo, mock, err := MockDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	type mockBehavior func(mock sqlmock.Sqlmock, returnedId int, book models.Book)
	tests := []struct {
		name         string
		inputBook    models.Book
		returnedId   int
		mockBehavior mockBehavior
		expectError  bool
	}{
		{
			name: "Ok",
			inputBook: models.Book{
				Name:   "hello",
				Price:  45.99,
				Genre:  1,
				Amount: 8,
			},
			returnedId: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, returnedId int, book models.Book) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO \"books\"").
					WithArgs(book.Name, book.Price, book.Genre, book.Amount).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(returnedId))
				mock.ExpectCommit()
			},
		},
		{
			name: "Empty field",
			inputBook: models.Book{
				Name:   "",
				Price:  45.99,
				Genre:  1,
				Amount: 8,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, returnedId int, book models.Book) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO \"books\"").
					WithArgs(book.Name, book.Price, book.Genre, book.Amount).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(returnedId).RowError(0, errors.New("insert error")))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(mock, test.returnedId, test.inputBook)
			id, err := repo.CreateBook(test.inputBook)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.returnedId, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteBookByID(t *testing.T) {
	repo, mock, err := MockDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	type mockBehavior func(inputId int)
	tests := []struct {
		name         string
		inputId      int
		mockBehavior mockBehavior
		expectError  bool
	}{
		{
			name:    "Ok",
			inputId: 3,
			mockBehavior: func(inputId int) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs(inputId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "Id not found",
			mockBehavior: func(inputId int) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE")).
					WithArgs(inputId).WillReturnError(errors.New("id not found"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.inputId)
			err := repo.DeleteBookByID(test.inputId)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetBooks(t *testing.T) {
	repo, mock, err := MockDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	type mockBehavior func(filterCondition map[string][]string)
	tests := []struct {
		name            string
		mockBehavior    mockBehavior
		filterCondition map[string][]string
		expectedBooks   []models.Book
		expectError     bool
	}{
		{
			name: "Ok",
			mockBehavior: func(filterCondition map[string][]string) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "genre", "amount"}).
						AddRow(1, "book1", 3.7, 1, 1).
						AddRow(2, "book2", 4.7, 2, 2).
						AddRow(3, "book3", 5.7, 3, 3))
			},
			expectedBooks: []models.Book{
				{ID: 1, Name: "book1", Price: 3.7, Genre: 1, Amount: 1},
				{ID: 2, Name: "book2", Price: 4.7, Genre: 2, Amount: 2},
				{ID: 3, Name: "book3", Price: 5.7, Genre: 3, Amount: 3},
			},
		},
		{
			name:            "Filter OK",
			filterCondition: map[string][]string{"genre": {"1"}},
			mockBehavior: func(filterCondition map[string][]string) {
				genreId := filterCondition["genre"][0]
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(0, genreId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "genre", "amount"}).
						AddRow(1, "book1", 3.7, 1, 1))
			},
			expectedBooks: []models.Book{
				{ID: 1, Name: "book1", Price: 3.7, Genre: 1, Amount: 1},
			},
		},
		{
			name:            "Filter returns empty array OK",
			filterCondition: map[string][]string{"genre": {"1"}},
			mockBehavior: func(filterCondition map[string][]string) {
				genreId := filterCondition["genre"][0]
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(0, genreId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "genre", "amount"}))
			},
			expectedBooks: []models.Book{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.filterCondition)
			books, err := repo.GetBooks(test.filterCondition)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBooks, books)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetBookByID(t *testing.T) {
	repo, mock, err := MockDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	type mockBehavior func(inputId int)
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		inputId      int
		expectedBook models.Book
		expectError  bool
	}{
		{
			name: "Ok",
			mockBehavior: func(inputId int) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(inputId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "genre", "amount"}).
						AddRow(inputId, "book1", 1.11, 2, 9))
			},
			inputId: 2,
			expectedBook: models.Book{
				ID:     2,
				Name:   "book1",
				Price:  1.11,
				Genre:  2,
				Amount: 9,
			},
		},
		{
			name: "Id not found",
			mockBehavior: func(inputId int) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(inputId).WillReturnError(errors.New("id not found"))
			},
			inputId:     2,
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.inputId)
			book, err := repo.GetBookByID(test.inputId)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBook, book)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateBookByID(t *testing.T) {
	repo, mock, err := MockDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	type mockBehavior func(inputId int, inputBook models.Book)
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		inputId      int
		inputBook    models.Book
		expectError  bool
	}{
		{
			name: "Ok",
			mockBehavior: func(inputId int, inputBook models.Book) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
					WithArgs(inputBook.Name, inputBook.Price, inputBook.Genre, inputBook.Amount, inputId, inputId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			inputId: 1,
			inputBook: models.Book{
				ID:     1,
				Name:   "book1",
				Price:  1.11,
				Genre:  2,
				Amount: 9,
			},
		},
		{
			name: "id not found",
			mockBehavior: func(inputId int, inputBook models.Book) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
					WithArgs(inputBook.Name, inputBook.Price, inputBook.Genre, inputBook.Amount, inputId, inputId).
					WillReturnError(errors.New("id not found"))
				mock.ExpectRollback()
			},
			inputId: 1,
			inputBook: models.Book{
				ID:     1,
				Name:   "book1",
				Price:  1.11,
				Genre:  2,
				Amount: 9,
			},
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.inputId, test.inputBook)
			err := repo.UpdateBookByID(test.inputId, test.inputBook)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
