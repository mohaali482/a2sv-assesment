package repository

import (
	"context"

	"github.com/mohaali482/a2sv-assesment/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	GetUsers(ctx context.Context) ([]*domain.User, error)
	Insert(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

const userCollection = "users"

type UserRepositoryImpl struct {
	database *mongo.Database
}

func NewUserRepository(database *mongo.Database) UserRepository {
	return &UserRepositoryImpl{database: database}
}

// Delete implements UserRepository.
func (u *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindByEmail implements UserRepository.
func (u *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	panic("unimplemented")
}

// FindByID implements UserRepository.
func (u *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.User, error) {
	panic("unimplemented")
}

// GetUsers implements UserRepository.
func (u *UserRepositoryImpl) GetUsers(ctx context.Context) ([]*domain.User, error) {
	panic("unimplemented")
}

// Insert implements UserRepository.
func (u *UserRepositoryImpl) Insert(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID().Hex()
	_, err := u.database.Collection(userCollection).InsertOne(ctx, user)
	return err
}

// Update implements UserRepository.
func (u *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "email", Value: user.Email},
			{Key: "password", Value: user.Password},
			{Key: "full_name", Value: user.FullName},
			{Key: "role", Value: user.Role},
		}},
	}

	_, err := u.database.Collection(userCollection).UpdateOne(ctx, filter, update)
	return err
}
