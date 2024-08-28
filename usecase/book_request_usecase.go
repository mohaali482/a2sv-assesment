package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/mohaali482/a2sv-assesment/domain"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
	"github.com/mohaali482/a2sv-assesment/repository"
)

type AddBorrowRequestRequest struct {
	BookID string `json:"book_id" validate:"required"`
}

type UpdateBorrowRequestRequest struct {
	Status string `json:"status" validate:"required"`
}

type BookRequestUseCase interface {
	AddBorrowRequest(ctx context.Context, userID string, form AddBorrowRequestRequest) (*domain.BorrowRequest, error)
	GetBorrowRequestByID(ctx context.Context, id string) (*domain.BorrowRequest, error)
	GetAllBorrowRequest(ctx context.Context, filter domain.BorrowRequestFilter) ([]*domain.BorrowRequest, error)
	ChangeBorrowStatus(ctx context.Context, id string, form UpdateBorrowRequestRequest) error
	DeleteBorrowRequest(ctx context.Context, id string) error
}

type BookRequestUseCaseImpl struct {
	repo repository.BookRequestRepository
}

func NewBookRequestUseCase(repo repository.BookRequestRepository) BookRequestUseCase {
	return &BookRequestUseCaseImpl{repo: repo}
}

// AddBorrowRequest implements BookRequestUseCase.
func (b *BookRequestUseCaseImpl) AddBorrowRequest(ctx context.Context, userID string, form AddBorrowRequestRequest) (*domain.BorrowRequest, error) {
	validate := validator.New()
	err := infrastructure.Validate(validate, form)
	if err != nil {
		return nil, err
	}

	exisitingRequest, _ := b.repo.GetBorrowRequestByUserIDAndBookID(ctx, userID, form.BookID)
	if exisitingRequest != nil {
		if exisitingRequest.Status == domain.BorrowRequestStatusPending {
			return nil, domain.ErrBorrowRequestAlreadyPending
		} else if exisitingRequest.Status == domain.BorrowRequestStatusApproved {
			return nil, domain.ErrBorrowRequestAlreadyApproved
		}
	}

	request := &domain.BorrowRequest{
		BookID: form.BookID,
		UserID: userID,
	}

	request.SetPending()
	err = b.repo.AddBorrowRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// ChangeBorrowStatus implements BookRequestUseCase.
func (b *BookRequestUseCaseImpl) ChangeBorrowStatus(ctx context.Context, id string, form UpdateBorrowRequestRequest) error {
	validate := validator.New()
	err := infrastructure.Validate(validate, form)
	if err != nil {
		return err
	}

	request, err := b.repo.GetBorrowRequestByID(ctx, id)
	if err != nil {
		return err
	}

	switch form.Status {
	case domain.BorrowRequestStatusApproved:
		request.SetApproved()
	case domain.BorrowRequestStatusRejected:
		request.SetRejected()
	default:
		return domain.ErrBorrowRequestStatusInvalid
	}

	return b.repo.UpdateBorrowRequest(ctx, request)
}

// DeleteBorrowRequest implements BookRequestUseCase.
func (b *BookRequestUseCaseImpl) DeleteBorrowRequest(ctx context.Context, id string) error {
	return b.repo.DeleteBorrowRequest(ctx, id)
}

// GetAllBorrowRequest implements BookRequestUseCase.
func (b *BookRequestUseCaseImpl) GetAllBorrowRequest(ctx context.Context, filter domain.BorrowRequestFilter) ([]*domain.BorrowRequest, error) {
	return b.repo.GetAllBorrowRequest(ctx, filter)
}

// GetBorrowRequestByID implements BookRequestUseCase.
func (b *BookRequestUseCaseImpl) GetBorrowRequestByID(ctx context.Context, id string) (*domain.BorrowRequest, error) {
	return b.repo.GetBorrowRequestByID(ctx, id)
}
