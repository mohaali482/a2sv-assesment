package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/controllers"
)

func NewLogRoute(r *gin.RouterGroup) {
	r.GET("/logs", controllers.LogController)
}
