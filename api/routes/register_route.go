package routes

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/api/controllers"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
	"github.com/mohaali482/a2sv-assesment/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRegisterRoute(db *mongo.Database, group *gin.RouterGroup) {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	es := infrastructure.NewEmailService(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)

	ur := repository.NewUserRepository(db)
	ps := infrastructure.NewPasswordService()
	vs := infrastructure.NewVerificationService()

	rc := controllers.NewRegisterController(
		usecase.NewRegisterUseCaseImpl(ur, ps, es, vs),
	)

	group.POST("/register", rc.Register)
}
