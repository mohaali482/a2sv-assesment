package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/controllers"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
	"github.com/mohaali482/a2sv-assesment/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewVerifyRoute(db *mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	vs := infrastructure.NewVerificationService()

	vc := controllers.NewVerifyController(
		usecase.NewVerifyUseCaseImpl(ur, vs),
	)

	group.POST("/verify-email", vc.Verify)
}
