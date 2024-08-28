package infrastructure

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/mohaali482/a2sv-assesment/domain"
)

type VerificationService interface {
	GenerateToken(user *domain.User) (string, error)
}

type VerificationServiceImpl struct {
}

func NewVerificationService() VerificationService {
	return &VerificationServiceImpl{}
}

func (v *VerificationServiceImpl) GenerateToken(user *domain.User) (string, error) {
	token := sha1.New()
	_, err := token.Write([]byte(user.Email + user.Password))
	if err != nil {
		return "", err
	}

	tokenString := string(token.Sum(nil))
	tokenString = base64.URLEncoding.EncodeToString([]byte(tokenString))
	return tokenString, nil
}
