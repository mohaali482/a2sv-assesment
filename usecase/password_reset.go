package usecase

import (
	"context"

	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
)

type PasswordResetForm struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordUpdateForm struct {
	ID          string `json:"id" validate:"required"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type PasswordResetUseCase interface {
	PasswordReset(ctx context.Context, form PasswordResetForm) error
	PasswordUpdate(ctx context.Context, form PasswordUpdateForm) error
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
	user, err := p.repo.FindByEmail(ctx, form.Email)
	if err != nil {
		return err
	}

	token, err := p.verificationService.GenerateToken(user)
	if err != nil {
		return err
	}

	passwordResetLink := "http://localhost:8080/users/password-update?id=" + user.ID + "&token=" + token

	err = p.emailService.Send(user.Email, "Password Reset", passwordResetLink)
	if err != nil {
		return err
	}

	return nil
}

// PasswordUpdate implements PasswordResetUseCase.
func (p *PasswordResetUseCaseImpl) PasswordUpdate(ctx context.Context, form PasswordUpdateForm) error {
	user, err := p.repo.FindByID(ctx, form.ID)
	if err != nil {
		return err
	}

	expectedToken, err := p.verificationService.GenerateToken(user)
	if err != nil {
		return err
	}

	if expectedToken != form.NewPassword {
		return domain.ErrInvalidVerificationCode
	}

	hashedPassword, err := p.passwordService.Hash(form.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	err = p.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
