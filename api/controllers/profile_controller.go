package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

type ProfileController struct {
	profileUseCase usecase.ProfileUsecase
}

func NewProfileController(profileUseCase usecase.ProfileUsecase) *ProfileController {
	return &ProfileController{profileUseCase}
}

func (p *ProfileController) Profile(c *gin.Context) {
	userID := c.Value("user_id").(string)

	user, err := p.profileUseCase.Profile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}
