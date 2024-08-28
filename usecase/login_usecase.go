package usecase

import (
	"context"

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

func NewLoginUseCaseImpl(repo repository.UserRepository) LoginUseCase {
	return &LoginUseCaseImpl{repo: repo}
}

// Login implements LoginUseCase.
func (l *LoginUseCaseImpl) Login(ctx context.Context, form LoginForm) (*domain.Token, error) {
	validate := validator.New()
	err := infrastructure.Validate(validate, form)
	if err != nil {
		return nil, err
	}

	user, err := l.repo.FindByEmail(ctx, form.Email)
	if !user.Verified {
		return nil, domain.ErrUserNotVerified
	}

	if !l.passwordService.Compare(user.Password, form.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := l.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return token, nil
}
