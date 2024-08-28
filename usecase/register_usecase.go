package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
)

type RegisterForm struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterUseCase interface {
	Register(ctx context.Context, form RegisterForm) error
	GenerateVerificationLink(ctx context.Context, user *domain.User) (string, error)
}

type RegisterUseCaseImpl struct {
	repo                repository.UserRepository
	passwordService     infrastructure.PasswordService
	emailService        infrastructure.EmailService
	verificationService infrastructure.VerificationService
}

func NewRegisterUseCaseImpl(repo repository.UserRepository, passwordService infrastructure.PasswordService) RegisterUseCase {
	return &RegisterUseCaseImpl{repo: repo, passwordService: passwordService}
}

// Register implements RegisterUseCase.
func (r *RegisterUseCaseImpl) Register(ctx context.Context, form RegisterForm) error {
	validate := validator.New()
	err := infrastructure.Validate(validate, form)
	if err != nil {
		return err
	}

	user := &domain.User{
		FullName: form.FullName,
		Email:    form.Email,
	}

	hashedPassword, err := r.passwordService.Hash(form.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Verified = false
	user.Role = domain.RoleUser

	err = r.repo.Insert(ctx, user)
	if err != nil {
		return err
	}

	verificationLink, err := r.GenerateVerificationLink(ctx, user)
	if err != nil {
		return err
	}

	err = r.emailService.Send(user.Email, "Email Verification", verificationLink)

	return err
}

// GenerateVerificationLink implements RegisterUseCase.
func (r *RegisterUseCaseImpl) GenerateVerificationLink(ctx context.Context, user *domain.User) (string, error) {
	tokenString, err := r.verificationService.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:8000/users/verify/%s/%s", user.ID, tokenString), nil
}
