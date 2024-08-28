package repository

import (
	"context"
	"log"

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
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := u.database.Collection(userCollection).DeleteOne(ctx, filter)
	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return err
}

// FindByEmail implements UserRepository.
func (u *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	var user domain.User
	err := u.database.Collection(userCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

// FindByID implements UserRepository.
func (u *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.User, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var user domain.User
	err := u.database.Collection(userCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil

}

// GetUsers implements UserRepository.
func (u *UserRepositoryImpl) GetUsers(ctx context.Context) ([]*domain.User, error) {
	cursor, err := u.database.Collection(userCollection).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var users []*domain.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
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
			{Key: "verified", Value: user.Verified},
		}},
	}

	result, err := u.database.Collection(userCollection).UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	if err != nil {
		log.Default().Println("Error updating user in repository:", err)
		return err
	}

	return nil
}
