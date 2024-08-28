package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

type PasswordController struct {
	passwordUseCase usecase.PasswordResetUseCase
}

func NewPasswordController(passwordUseCase usecase.PasswordResetUseCase) *PasswordController {
	return &PasswordController{passwordUseCase}
}

func (p *PasswordController) PasswordReset(c *gin.Context) {
	var form usecase.PasswordResetForm
	if err := c.ShouldBindJSON(&form); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := p.passwordUseCase.PasswordReset(c.Request.Context(), form)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, infrastructure.ReturnErrorResponse(err))
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "password reset link sent"})
}

func (p *PasswordController) PasswordUpdate(c *gin.Context) {
	var form usecase.PasswordUpdateForm
	if err := c.ShouldBindJSON(&form); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Query("id")
	token := c.Query("token")

	err := p.passwordUseCase.PasswordUpdate(c.Request.Context(), id, token, form)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, infrastructure.ReturnErrorResponse(err))
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "password updated successfully"})
}
