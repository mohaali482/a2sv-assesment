package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/controllers"
	"github.com/mohaali482/a2sv-assesment/repository"
	"github.com/mohaali482/a2sv-assesment/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRoute(db *mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	uc := controllers.NewUserController(
		usecase.NewUserUsecaseImpl(ur),
	)

	group.POST("/users", uc.GetUsers)
	group.DELETE("/users/:id", uc.DeleteUser)
}
