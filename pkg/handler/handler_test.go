package handler

import (
	"bytes"
	"errors"
	"fmt"
	restapi "github.com/TenderLimbo/rest-api"
	"github.com/TenderLimbo/rest-api/pkg/service"
	mock_service "github.com/TenderLimbo/rest-api/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandler_CreateBook(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBooksManager, book restapi.Book)
	tests := []struct {
		name                 string
		inputBody            string
		inputBook            restapi.Book
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name": "Book1", "price": 0, "genre": 1, "amount": 0}`,
			inputBook: restapi.Book{
				Model:  restapi.Model{},
				Name:   "Book1",
				Price:  0,
				Genre:  1,
				Amount: 0,
			},
			mockBehavior: func(r *mock_service.MockBooksManager, book restapi.Book) {
				r.EXPECT().CreateBook(book).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Some fields missing",
			inputBody:            `{"price": 67.88, "genre": 1, "amount": 5}`,
			mockBehavior:         func(r *mock_service.MockBooksManager, book restapi.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid input"}`,
		},
		{
			name:                 "Invalid genre",
			inputBody:            `{"name": "hello", "price": 67.88, "genre": 6, "amount": 7}`,
			mockBehavior:         func(r *mock_service.MockBooksManager, book restapi.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid input"}`,
		},
		{
			name:      "Not unique name",
			inputBody: `{"name": "hello", "price": 67.88, "genre": 1, "amount": 7}`,
			inputBook: restapi.Book{
				Model:  restapi.Model{},
				Name:   "hello",
				Price:  67.88,
				Genre:  1,
				Amount: 7,
			},
			mockBehavior: func(r *mock_service.MockBooksManager, book restapi.Book) {
				r.EXPECT().CreateBook(book).Return(0, errors.New("input book name isn't unique"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"input book name isn't unique"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock_service.NewMockBooksManager(c)
			test.mockBehavior(mockManager, test.inputBook)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := gin.New()
			r.POST("/books", handler.CreateBook)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/books",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}

}

func TestHandler_DeleteBookByID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBooksManager, id interface{})
	tests := []struct {
		name                 string
		inputId              interface{}
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Id invalid",
			inputId:              "knekndijf",
			mockBehavior:         func(r *mock_service.MockBooksManager, id interface{}) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:    "Id OK",
			inputId: 1,
			mockBehavior: func(r *mock_service.MockBooksManager, id interface{}) {
				r.EXPECT().DeleteBookByID(id).Return(nil)
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: ``,
		},
		{
			name:    "Id not found",
			inputId: 1,
			mockBehavior: func(r *mock_service.MockBooksManager, id interface{}) {
				r.EXPECT().DeleteBookByID(id).Return(errors.New("id not found"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"id not found"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock_service.NewMockBooksManager(c)
			test.mockBehavior(mockManager, test.inputId)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := gin.New()
			r.DELETE("/books/:id", handler.DeleteBookByID)
			target := fmt.Sprintf("/books/%v", test.inputId)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", target, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetBookByID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBooksManager, id interface{})
	tests := []struct {
		name                 string
		inputId              interface{}
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Id invalid",
			inputId:              "knekndijf",
			mockBehavior:         func(r *mock_service.MockBooksManager, id interface{}) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:    "Id OK",
			inputId: 1,
			mockBehavior: func(r *mock_service.MockBooksManager, id interface{}) {
				r.EXPECT().GetBookByID(id).Return(restapi.Book{
					Model:  restapi.Model{ID: 1},
					Name:   "hello",
					Price:  4.32,
					Genre:  2,
					Amount: 9,
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"name":"hello","price":4.32,"genre":2,"amount":9}`,
		},
		{
			name:    "Id not found",
			inputId: 1,
			mockBehavior: func(r *mock_service.MockBooksManager, id interface{}) {
				r.EXPECT().GetBookByID(id).Return(restapi.Book{}, errors.New("id not found"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"id not found"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock_service.NewMockBooksManager(c)
			test.mockBehavior(mockManager, test.inputId)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := gin.New()
			r.GET("/books/:id", handler.GetBookByID)
			target := fmt.Sprintf("/books/%v", test.inputId)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", target, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetBooks(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBooksManager, filterCondition map[string][]string)
	tests := []struct {
		name                 string
		filterCondition      map[string][]string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Invalid filter condition",
			filterCondition:      map[string][]string{"kjknmlm": {"7"}},
			mockBehavior:         func(r *mock_service.MockBooksManager, filterCondition map[string][]string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid filter condition"}`,
		},
		{
			name:                 "Invalid genre id in filter condition",
			filterCondition:      map[string][]string{"genre": {"0"}},
			mockBehavior:         func(r *mock_service.MockBooksManager, filterCondition map[string][]string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid filter condition"}`,
		},
		{
			name:            "Get All Ok",
			filterCondition: map[string][]string{},
			mockBehavior: func(r *mock_service.MockBooksManager, filterCondition map[string][]string) {
				r.EXPECT().GetBooks(filterCondition).Return([]restapi.Book{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock_service.NewMockBooksManager(c)
			test.mockBehavior(mockManager, test.filterCondition)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := gin.New()
			r.GET("/books", handler.GetBooks)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/books", nil)
			values := url.Values(test.filterCondition)
			req.URL.RawQuery = values.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_UpdateBookByID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBooksManager, id interface{}, book restapi.Book)
	tests := []struct {
		name                 string
		inputBody            string
		inputId              interface{}
		inputBook            restapi.Book
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Update invalid id",
			inputId:              78.99,
			mockBehavior:         func(r *mock_service.MockBooksManager, id interface{}, book restapi.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid id"}`,
		},
		{
			name:                 "Update invalid input",
			inputId:              1,
			inputBody:            `{"name": "Book1", "price": -78, "genre": 1, "amount": 0}`,
			mockBehavior:         func(r *mock_service.MockBooksManager, id interface{}, book restapi.Book) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid input"}`,
		},
		{
			name:      "Update id not found",
			inputId:   1,
			inputBody: `{"name": "Book1", "price": 0, "genre": 1, "amount": 0}`,
			inputBook: restapi.Book{
				Model:  restapi.Model{},
				Name:   "Book1",
				Price:  0,
				Genre:  1,
				Amount: 0,
			},
			mockBehavior: func(r *mock_service.MockBooksManager, id interface{}, book restapi.Book) {
				r.EXPECT().UpdateBookByID(id, book).Return(errors.New("id not found"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"id not found"}`,
		},
		{
			name:      "Update ok",
			inputId:   1,
			inputBody: `{"name": "Book1", "price": 0, "genre": 1, "amount": 0}`,
			inputBook: restapi.Book{
				Model:  restapi.Model{},
				Name:   "Book1",
				Price:  0,
				Genre:  1,
				Amount: 0,
			},
			mockBehavior: func(r *mock_service.MockBooksManager, id interface{}, book restapi.Book) {
				r.EXPECT().UpdateBookByID(id, book).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"name":"Book1","price":0,"genre":1,"amount":0}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockManager := mock_service.NewMockBooksManager(c)
			test.mockBehavior(mockManager, test.inputId, test.inputBook)

			services := service.NewService(mockManager)
			handler := Handler{services}

			r := gin.New()
			r.PUT("/books/:id", handler.UpdateBookByID)
			target := fmt.Sprintf("/books/%v", test.inputId)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", target, bytes.NewBufferString(test.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
