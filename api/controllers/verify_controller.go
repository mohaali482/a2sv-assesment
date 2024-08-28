package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

type VerifyController struct {
	verifyUseCase usecase.VerifyUseCase
}

func NewVerifyController(verifyUseCase usecase.VerifyUseCase) *VerifyController {
	return &VerifyController{verifyUseCase}
}

func (v *VerifyController) Verify(c *gin.Context) {
	userID := c.Param("id")
	token := c.Query("token")

	err := v.verifyUseCase.Verify(c.Request.Context(), userID, token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user verified successfully"})
}
