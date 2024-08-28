package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/middleware"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(db *mongo.Database, gin *gin.Engine) {
	js := infrastructure.NewJWTService(os.Getenv("JWT_SECRET"))
	middleware := middleware.NewMiddlewareImpl(js)

	publicRoutes := gin.Group("users")

	NewRegisterRoute(db, publicRoutes)
	NewVerifyRoute(db, publicRoutes)
	NewLoginRoute(db, publicRoutes)
	NewRefreshRoute(publicRoutes)
	NewPasswordRoute(db, publicRoutes)

	privateRoutes := gin.Group("users")
	privateRoutes.Use(middleware.AuthMiddleware())
	NewProfileRoute(db, privateRoutes)

	adminRoutes := gin.Group("admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(middleware.AdminMiddleware())
	NewUserRoute(db, adminRoutes)

}
