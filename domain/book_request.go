package domain

import "errors"

var ErrBookNotFound = errors.New("book not found")
var ErrBookRequestNotFound = errors.New("book request not found")
var ErrBorrowRequestAlreadyPending = errors.New("borrow request already pending")
var ErrBorrowRequestAlreadyApproved = errors.New("borrow request already approved")
var ErrBorrowRequestStatusInvalid = errors.New("invalid borrow request status")

var BorrowRequestStatusPending = "pending"
var BorrowRequestStatusApproved = "approved"
var BorrowRequestStatusRejected = "rejected"

type BorrowRequestFilter struct {
	Status string // pending, approved, rejected default: all
	Order  string // asc, desc default: asc
}

type BorrowRequest struct {
	ID     string `json:"id" bson:"_id"`
	BookID string `json:"book_id" bson:"book_id"`
	UserID string `json:"user_id" bson:"user_id"`
	Status string `json:"status" bson:"status"`
}

func (b *BorrowRequest) SetApproved() {
	b.Status = BorrowRequestStatusApproved
}

func (b *BorrowRequest) SetRejected() {
	b.Status = BorrowRequestStatusRejected
}

func (b *BorrowRequest) SetPending() {
	b.Status = BorrowRequestStatusPending
}
