package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

type RefreshForm struct {
	Refresh string `json:"refresh"`
}

type RefreshController struct {
	refreshUseCase usecase.RefreshUseCase
}

func NewRefreshController(refreshUseCase usecase.RefreshUseCase) *RefreshController {
	return &RefreshController{refreshUseCase}
}

func (r *RefreshController) Refresh(c *gin.Context) {
	var form RefreshForm
	if err := c.ShouldBindJSON(&form); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.refreshUseCase.Refresh(c.Request.Context(), form.Refresh)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, infrastructure.ReturnErrorResponse(err))
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access": token,
	})
}
