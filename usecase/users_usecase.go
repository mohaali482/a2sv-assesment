package usecase

import (
	"context"

	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/repository"
)

type UserUsecase interface {
	GetUsers(ctx context.Context) ([]*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUserUsecaseImpl(repo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{
		repo: repo,
	}
}

func (u *UserUsecaseImpl) GetUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserUsecaseImpl) DeleteUser(ctx context.Context, id string) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
