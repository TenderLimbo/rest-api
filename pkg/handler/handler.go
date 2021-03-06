package handler

import (
	"github.com/TenderLimbo/rest-api/models"
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
	if len(filterCondition) != 0 {
		if !filterCondition.Has("genre") {
			NewErrorResponse(ctx, http.StatusBadRequest, "invalid filter condition")
			return
		}
		genreID, err := strconv.Atoi(filterCondition.Get("genre"))
		if err != nil || genreID < 1 || genreID > 3 {
			NewErrorResponse(ctx, http.StatusBadRequest, "invalid filter condition")
			return
		}
	}
	books, err := h.service.GetBooks(filterCondition)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (h *Handler) GetBookByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}
	book, err := h.service.GetBookByID(id)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (h *Handler) CreateBook(ctx *gin.Context) {
	var newBook models.Book
	if err := ctx.BindJSON(&newBook); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid input")
		return
	}
	id, err := h.service.CreateBook(newBook)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) DeleteBookByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}
	if err = h.service.DeleteBookByID(id); err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, StatusResponse{"ok"})
}

func (h *Handler) UpdateBookByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid id")
		return
	}
	var newBook models.Book
	if err = ctx.BindJSON(&newBook); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid input")
		return
	}
	if err = h.service.UpdateBookByID(id, newBook); err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	newBook.ID = id
	ctx.JSON(http.StatusOK, newBook)
}
