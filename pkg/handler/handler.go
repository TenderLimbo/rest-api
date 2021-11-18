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
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (h *Handler) GetBookByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
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
	var newBook restapi.Book
	if err := ctx.BindJSON(&newBook); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := restapi.ValidateBook(newBook); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.CreateBook(newBook)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) DeleteBookByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err = h.service.DeleteBookByID(id); err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}

func (h *Handler) UpdateBookByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var newBook restapi.Book
	if err = ctx.BindJSON(&newBook); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := restapi.ValidateBook(newBook); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err = h.service.UpdateBookByID(id, newBook); err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, StatusResponse{"ok"})
}
