package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
)

type PasswordResetForm struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordUpdateForm struct {
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type PasswordResetUseCase interface {
	PasswordReset(ctx context.Context, form PasswordResetForm) error
	PasswordUpdate(ctx context.Context, id, token string, form PasswordUpdateForm) error
}

type PasswordResetUseCaseImpl struct {
	repo                repository.UserRepository
	passwordService     infrastructure.PasswordService
	emailService        infrastructure.EmailService
	verificationService infrastructure.VerificationService
}

func NewPasswordResetUseCaseImpl(
	repo repository.UserRepository,
	passwordService infrastructure.PasswordService,
	emailService infrastructure.EmailService,
	verificationService infrastructure.VerificationService,
) PasswordResetUseCase {
	return &PasswordResetUseCaseImpl{
		repo:                repo,
		passwordService:     passwordService,
		emailService:        emailService,
		verificationService: verificationService,
	}
}

// PasswordReset implements PasswordResetUseCase.
func (p *PasswordResetUseCaseImpl) PasswordReset(ctx context.Context, form PasswordResetForm) error {
	validate := validator.New()
	err := infrastructure.Validate(validate, form)
	if err != nil {
		return err
	}

	user, err := p.repo.FindByEmail(ctx, form.Email)
	if err != nil {
		return err
	}

	token, err := p.verificationService.GenerateToken(user)
	if err != nil {
		return err
	}

	passwordResetLink := "http://localhost:8000/users/password-update?id=" + user.ID + "&token=" + token

	err = p.emailService.Send(user.Email, "Password Reset", passwordResetLink)
	if err != nil {
		return err
	}

	log.Default().Println("User with the email of", user.Email, "requested for reset password")
	return nil
}

// PasswordUpdate implements PasswordResetUseCase.
func (p *PasswordResetUseCaseImpl) PasswordUpdate(ctx context.Context, id, token string, form PasswordUpdateForm) error {
	validator := validator.New()
	err := infrastructure.Validate(validator, form)
	if err != nil {
		return err
	}

	if form.NewPassword != form.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	user, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	expectedToken, err := p.verificationService.GenerateToken(user)
	if err != nil {
		return err
	}

	if expectedToken != token {
		return domain.ErrInvalidVerificationCode
	}

	hashedPassword, err := p.passwordService.Hash(form.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	err = p.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	log.Default().Println("User with the email of", user.Email, "reseted password successfully")
	return nil
}
