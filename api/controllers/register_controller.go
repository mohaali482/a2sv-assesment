package controllers

import (
	"github.com/gin-gonic/gin"
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := r.registerUseCase.Register(c.Request.Context(), form)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user registered successfully"})
}
