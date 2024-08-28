package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/controllers"
	"github.com/mohaali482/a2sv-assesment/repository"
	"github.com/mohaali482/a2sv-assesment/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewProfileRoute(db *mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	pc := controllers.NewProfileController(
		usecase.NewProfileUsecaseImpl(ur),
	)

	group.GET("/profile", pc.Profile)
}
