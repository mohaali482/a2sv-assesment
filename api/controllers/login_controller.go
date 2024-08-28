package controllers

import (
	"github.com/gin-gonic/gin"
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := l.loginUsecase.Login(c.Request.Context(), form)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, token)
}
