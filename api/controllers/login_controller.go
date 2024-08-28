package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

type LoginController struct {
	loginUsecase usecase.LoginUseCase
}

func NewLoginController(loginUsecase usecase.LoginUseCase) *LoginController {
	return &LoginController{loginUsecase}
}

func (l *LoginController) Login(c *gin.Context) {
	var form usecase.LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := l.loginUsecase.Login(c.Request.Context(), form)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, infrastructure.ReturnErrorResponse(err))
			return
		}

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, token)
}
