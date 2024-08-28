package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/controllers"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/usecase"
)

func NewRefreshRoute(group *gin.RouterGroup) {
	js := infrastructure.NewJWTService(os.Getenv("JWT_SECRET"))

	rc := controllers.NewRefreshController(
		usecase.NewRefreshUseCaseImpl(js),
	)

	group.POST("/refresh", rc.Refresh)
}
