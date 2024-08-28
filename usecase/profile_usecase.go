package usecase

import (
	"context"

	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/repository"
)

type ProfileUsecase interface {
	Profile(ctx context.Context, id string) (*domain.User, error)
}

type ProfileUsecaseImpl struct {
	repo repository.UserRepository
}

func NewProfileUsecaseImpl(repo repository.UserRepository) ProfileUsecase {
	return &ProfileUsecaseImpl{
		repo: repo,
	}
}

func (p *ProfileUsecaseImpl) Profile(ctx context.Context, id string) (*domain.User, error) {
	user, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
