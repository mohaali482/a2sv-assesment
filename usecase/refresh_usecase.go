package usecase

import (
	"context"

	"github.com/mohaali482/a2sv-assesment/infrastructure"
)

type RefreshUseCase interface {
	Refresh(ctx context.Context, token string) (string, error)
}

type RefreshUseCaseImpl struct {
	jwtService infrastructure.JWTService
}

func NewRefreshUseCaseImpl(jwtService infrastructure.JWTService) RefreshUseCase {
	return &RefreshUseCaseImpl{jwtService: jwtService}
}

// Refresh implements RefreshUseCase.
func (r *RefreshUseCaseImpl) Refresh(ctx context.Context, token string) (string, error) {
	claims, err := r.jwtService.ValidateToken(token)
	tokenType, ok := claims["type"].(string)

	if !ok || tokenType != "refresh" {
		return "", infrastructure.ErrInvalidToken
	}

	if err != nil {
		return "", err
	}

	access, err := r.jwtService.GenerateAccess(claims)
	if err != nil {
		return "", err
	}

	return access, nil
}
