package usecase

import (
	"context"

	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
)

type VerifyUseCase interface {
	Verify(ctx context.Context, id, token string) error
}

type VerifyUseCaseImpl struct {
	repo                repository.UserRepository
	verificationService infrastructure.VerificationService
}

func NewVerifyUseCaseImpl(
	repo repository.UserRepository,
	verificationService infrastructure.VerificationService,
) VerifyUseCase {
	return &VerifyUseCaseImpl{repo: repo, verificationService: verificationService}
}

// Verify implements VerifyUseCase.
func (v *VerifyUseCaseImpl) Verify(ctx context.Context, id, token string) error {
	user, err := v.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	expectedToken, err := v.verificationService.GenerateToken(user)
	if err != nil {
		return err
	}

	if expectedToken != token {
		return domain.ErrInvalidVerificationCode
	}

	user.Verified = true
	err = v.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
