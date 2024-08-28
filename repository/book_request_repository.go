package repository

import (
	"context"
	"log"

	"github.com/mohaali482/a2sv-assesment/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookRequestCollection = "books-request"

type BookRequestRepository interface {
	AddBorrowRequest(ctx context.Context, request *domain.BorrowRequest) error
	UpdateBorrowRequest(ctx context.Context, request *domain.BorrowRequest) error
	GetBorrowRequestByUserIDAndBookID(ctx context.Context, userID, bookID string) (*domain.BorrowRequest, error)
	GetAllBorrowRequest(ctx context.Context, filter domain.BorrowRequestFilter) ([]*domain.BorrowRequest, error)
	GetBorrowRequestByID(ctx context.Context, id string) (*domain.BorrowRequest, error)
	DeleteBorrowRequest(ctx context.Context, id string) error
}

type BookRequestRepositoryImpl struct {
	db *mongo.Database
}

func NewBookRequestRepository(db *mongo.Database) BookRequestRepository {
	return &BookRequestRepositoryImpl{db: db}
}

func (r *BookRequestRepositoryImpl) AddBorrowRequest(ctx context.Context, request *domain.BorrowRequest) error {
	request.ID = primitive.NewObjectID().Hex()

	_, err := r.db.Collection(bookRequestCollection).InsertOne(ctx, request)
	if err != nil {
		log.Default().Println("Error adding borrow request in repository:", err)
		return err
	}

	return nil
}

func (r *BookRequestRepositoryImpl) UpdateBorrowRequest(ctx context.Context, request *domain.BorrowRequest) error {
	filter := bson.D{{Key: "_id", Value: request.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "book_id", Value: request.BookID},
			{Key: "user_id", Value: request.UserID},
			{Key: "status", Value: request.Status},
		}},
	}
	result, err := r.db.Collection(bookRequestCollection).UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		return domain.ErrBookRequestNotFound
	}

	if err != nil {
		log.Default().Println("Error updating borrow request in repository:", err)
		return err
	}

	return err
}

func (r *BookRequestRepositoryImpl) GetBorrowRequestByUserIDAndBookID(ctx context.Context, userID, bookID string) (*domain.BorrowRequest, error) {
	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "book_id", Value: bookID},
	}
	var request domain.BorrowRequest
	err := r.db.Collection(bookRequestCollection).FindOne(ctx, filter).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrBookRequestNotFound
		}

		log.Default().Println("Error getting borrow request in repository:", err)
		return nil, err
	}

	return &request, nil
}

func (r *BookRequestRepositoryImpl) GetAllBorrowRequest(ctx context.Context, filter domain.BorrowRequestFilter) ([]*domain.BorrowRequest, error) {
	query := bson.D{}
	if filter.Status == domain.BorrowRequestStatusPending || filter.Status == domain.BorrowRequestStatusApproved || filter.Status == domain.BorrowRequestStatusRejected {
		query = append(query, bson.E{Key: "status", Value: filter.Status})
	}

	statusSortOrder := bson.M{
		"pending":  1, // Pending first
		"approved": 2, // Approved second
		"rejected": 3, // Rejected third
	}

	// Pipeline to add a computed field for sorting
	pipeline := mongo.Pipeline{
		// Match stage to apply any filters
		bson.D{{Key: "$match", Value: query}},
		// Add a field 'statusOrder' based on the status
		bson.D{{Key: "$addFields", Value: bson.M{
			"statusOrder": bson.M{
				"$switch": bson.M{
					"branches": []bson.M{
						{"case": bson.M{"$eq": []interface{}{"$status", "pending"}}, "then": statusSortOrder["pending"]},
						{"case": bson.M{"$eq": []interface{}{"$status", "approved"}}, "then": statusSortOrder["approved"]},
						{"case": bson.M{"$eq": []interface{}{"$status", "rejected"}}, "then": statusSortOrder["rejected"]},
					},
					"default": 4, // Default to the lowest priority
				},
			},
		}}},
		// Sort by the computed 'statusOrder' field
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "statusOrder", Value: 1}, // Ascending order by statusOrder
		}}},
	}

	// Override the sort order if 'Order' is explicitly set
	if filter.Order == "asc" {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{
			{Key: "statusOrder", Value: 1}, // Ascending order
		}}})
	} else if filter.Order == "desc" {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{
			{Key: "statusOrder", Value: -1}, // Descending order
		}}})
	}

	cursor, err := r.db.Collection(bookRequestCollection).Aggregate(ctx, pipeline)
	if err != nil {
		log.Default().Println("Error getting all borrow request in repository:", err)
		return nil, err
	}

	var requests []*domain.BorrowRequest = make([]*domain.BorrowRequest, 0)
	if err = cursor.All(ctx, &requests); err != nil {
		log.Default().Println("Error decoding all borrow request in repository:", err)
		return nil, err
	}

	return requests, nil
}

func (r *BookRequestRepositoryImpl) GetBorrowRequestByID(ctx context.Context, id string) (*domain.BorrowRequest, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var request domain.BorrowRequest
	err := r.db.Collection(bookRequestCollection).FindOne(ctx, filter).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrBookRequestNotFound
		}

		log.Default().Println("Error getting borrow request in repository:", err)
		return nil, err
	}

	return &request, nil
}

func (r *BookRequestRepositoryImpl) DeleteBorrowRequest(ctx context.Context, id string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := r.db.Collection(bookRequestCollection).DeleteOne(ctx, filter)
	if result.DeletedCount == 0 {
		return domain.ErrBookRequestNotFound
	}

	if err != nil {
		log.Default().Println("Error deleting borrow request in repository:", err)
		return err
	}

	return nil
}
