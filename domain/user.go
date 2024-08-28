package domain

import "errors"

const (
	RoleAdmin = iota
	RoleUser
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserNotVerified = errors.New("user not verified")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidVerificationCode = errors.New("invalid verification code")

type User struct {
	ID       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	FullName string `json:"full_name" bson:"full_name"`
	Role     int    `json:"role" bson:"role"`
	Verified bool   `json:"verified" bson:"verified"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
