package usecase

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
)

type LoginForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUseCase interface {
	Login(ctx context.Context, form LoginForm) (*domain.Token, error)
}

type LoginUseCaseImpl struct {
	repo            repository.UserRepository
	passwordService infrastructure.PasswordService
	jwtService      infrastructure.JWTService
}

func NewLoginUseCaseImpl(repo repository.UserRepository, passwordService infrastructure.PasswordService, jwtService infrastructure.JWTService) LoginUseCase {
	return &LoginUseCaseImpl{
		repo:            repo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

// Login implements LoginUseCase.
func (l *LoginUseCaseImpl) Login(ctx context.Context, form LoginForm) (*domain.Token, error) {
	validate := validator.New()
	err := infrastructure.Validate(validate, form)
	if err != nil {
		log.Default().Println("Invalid form:", err)
		return nil, err
	}

	user, err := l.repo.FindByEmail(ctx, form.Email)
	if err != nil {
		log.Default().Println("Failed to find user by email:", err)
		return nil, err
	}

	if !user.Verified {
		log.Default().Println("User is not verified. Email:", user.Email)
		return nil, domain.ErrUserNotVerified
	}

	if !l.passwordService.Compare(user.Password, form.Password) {
		log.Default().Println("Invalid credentials from user with email of", user.Email)
		return nil, domain.ErrInvalidCredentials
	}

	token, err := l.jwtService.GenerateToken(user)
	if err != nil {
		log.Default().Println("Failed to generate token:", err)
		return nil, err
	}

	log.Default().Println("User with email:", user.Email, "successfully logged in")
	return token, nil
}
