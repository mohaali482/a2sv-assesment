package controllers

import (
	"github.com/gin-gonic/gin"
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := r.refreshUseCase.Refresh(c.Request.Context(), form.Refresh)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, token)
}
