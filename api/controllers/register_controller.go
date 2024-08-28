package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

type RegisterController struct {
	registerUseCase usecase.RegisterUseCase
}

func NewRegisterController(registerUseCase usecase.RegisterUseCase) *RegisterController {
	return &RegisterController{registerUseCase}
}

func (r *RegisterController) Register(c *gin.Context) {
	var form usecase.RegisterForm
	if err := c.ShouldBindJSON(&form); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.registerUseCase.Register(c.Request.Context(), form)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, infrastructure.ReturnErrorResponse(err))
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user registered successfully"})
}
