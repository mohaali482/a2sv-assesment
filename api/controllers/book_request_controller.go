package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

type BookRequestController struct {
	bookRequestUseCase usecase.BookRequestUseCase
}

func NewBookRequestController(bookRequestUseCase usecase.BookRequestUseCase) *BookRequestController {
	return &BookRequestController{bookRequestUseCase: bookRequestUseCase}
}

func (b *BookRequestController) AddBorrowRequest(c *gin.Context) {
	var form usecase.AddBorrowRequestRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Value("user_id").(string)

	request, err := b.bookRequestUseCase.AddBorrowRequest(c.Request.Context(), userID, form)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, infrastructure.ReturnErrorResponse(err))
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, request)
}

func (b *BookRequestController) GetBorrowRequestByID(c *gin.Context) {
	id := c.Param("id")
	request, err := b.bookRequestUseCase.GetBorrowRequestByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrBookRequestNotFound {
			c.JSON(404, gin.H{"error": "Borrow request not found"})
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, request)
}

func (b *BookRequestController) GetAllBorrowRequest(c *gin.Context) {
	filter := domain.BorrowRequestFilter{
		Status: c.Query("status"),
		Order:  c.Query("order"),
	}

	requests, err := b.bookRequestUseCase.GetAllBorrowRequest(c.Request.Context(), filter)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, requests)
}

func (b *BookRequestController) UpdateBorrowRequest(c *gin.Context) {
	id := c.Param("id")
	var form usecase.UpdateBorrowRequestRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := b.bookRequestUseCase.ChangeBorrowStatus(c.Request.Context(), id, form)

	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, infrastructure.ReturnErrorResponse(err))
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Borrow request status updated successfully"})
}

func (b *BookRequestController) DeleteBorrowRequest(c *gin.Context) {
	id := c.Param("id")
	err := b.bookRequestUseCase.DeleteBorrowRequest(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrBookRequestNotFound {
			c.JSON(404, gin.H{"error": "Borrow request not found"})
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Borrow request deleted successfully"})
}
