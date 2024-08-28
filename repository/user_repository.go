package repository

import (
	"context"

	"github.com/mohaali482/a2sv-assesment/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	GetUsers(ctx context.Context) ([]*domain.User, error)
	Insert(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}
