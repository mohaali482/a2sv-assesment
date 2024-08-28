package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/controllers"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
	"github.com/mohaali482/a2sv-assesment/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewLoginRoute(db *mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	ps := infrastructure.NewPasswordService()
	js := infrastructure.NewJWTService(os.Getenv("JWT_SECRET"))
	lc := controllers.NewLoginController(
		usecase.NewLoginUseCaseImpl(ur, ps, js),
	)

	group.POST("/login", lc.Login)
}
