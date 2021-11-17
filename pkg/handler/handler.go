package handler

import (
	restapi "github.com/TenderLimbo/rest-api"
	"github.com/TenderLimbo/rest-api/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	service service.BooksManager
}

func NewHandler(service service.BooksManager) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/books", h.GetBooks)
	router.GET("/books/:id", h.GetBookByID)
	router.POST("/books", h.CreateBook)
	router.DELETE("books/:id", h.DeleteBookByID)
	router.PUT("books/:id", h.UpdateBookByID)
	return router
}

func (h *Handler) GetBooks(ctx *gin.Context) {
	filterCondition := ctx.Request.URL.Query()
	books, err := h.service.GetBooks(filterCondition)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, books)
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (h *Handler) GetBookByID(ctx *gin.Context) {
	id := ctx.Param("id")
	book, err := h.service.GetBookByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "book " + id + " not found"})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (h *Handler) CreateBook(ctx *gin.Context) {
	var newBook restapi.Book
	if err := ctx.BindJSON(&newBook); err != nil || !restapi.ValidateBook(newBook) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "can't create book"})
		return
	}
	id, err := h.service.CreateBook(newBook)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "book named " + newBook.Name + " already exists"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "book " + strconv.Itoa(id) + " created"})
}

func (h *Handler) DeleteBookByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.service.DeleteBookByID(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "book " + id + " not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "book " + id + " is deleted"})
}

func (h *Handler) UpdateBookByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var newBook restapi.Book
	if err := ctx.BindJSON(&newBook); err != nil || !restapi.ValidateBook(newBook) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "can't update book"})
		return
	}
	if err := h.service.UpdateBookByID(id, newBook); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "book " + id + " not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "book " + id + " updated"})
}
