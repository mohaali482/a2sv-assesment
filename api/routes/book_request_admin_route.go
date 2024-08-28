package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/controllers"
	"github.com/mohaali482/a2sv-assesment/repository"
	"github.com/mohaali482/a2sv-assesment/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewBookRequestAdminRoute(db *mongo.Database, group *gin.RouterGroup) {
	bookRequestRepository := repository.NewBookRequestRepository(db)
	bookRequestUseCase := usecase.NewBookRequestUseCase(bookRequestRepository)
	bookRequestController := controllers.NewBookRequestController(bookRequestUseCase)

	group.GET("/", bookRequestController.GetAllBorrowRequest)
	group.PATCH("/:id/status", bookRequestController.UpdateBorrowRequest)
	group.DELETE("/:id", bookRequestController.DeleteBorrowRequest)
}
