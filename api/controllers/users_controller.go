package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

type UserController struct {
	userUseCase usecase.UserUsecase
}

func NewUserController(userUseCase usecase.UserUsecase) *UserController {
	return &UserController{userUseCase}
}

func (u *UserController) GetUsers(c *gin.Context) {
	users, err := u.userUseCase.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, users)
}

func (u *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := u.userUseCase.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user deleted successfully"})
}
