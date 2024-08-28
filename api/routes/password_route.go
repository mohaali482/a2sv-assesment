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

func NewPasswordRoute(db *mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	ps := infrastructure.NewPasswordService()
	vs := infrastructure.NewVerificationService()

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

	pc := controllers.NewPasswordController(
		usecase.NewPasswordResetUseCaseImpl(ur, ps, es, vs),
	)

	group.POST("/password-reset", pc.PasswordReset)
	group.POST("/password-update", pc.PasswordUpdate)
}
