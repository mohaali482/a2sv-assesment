package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mohaali482/a2sv-assesment/domain"
)

var ErrInvalidToken = errors.New("invalid token")

const AccessTokenDuration = 15 * time.Minute
const RefreshTokenDuration = 7 * 24 * time.Hour

type JWTService interface {
	GenerateToken(user *domain.User) (*domain.Token, error)
	GenerateAccess(claims jwt.MapClaims) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}

type JWTServiceImpl struct {
	secretKey string
}

func NewJWTService(secretKey string) JWTService {
	return &JWTServiceImpl{secretKey: secretKey}
}

// GenerateToken implements JWTService.
func (j *JWTServiceImpl) GenerateToken(user *domain.User) (*domain.Token, error) {
	accessTokenclaims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"type":  "access",
		"exp":   jwt.TimeFunc().Add(AccessTokenDuration).Unix(),
	}

	if user.IsAdmin() {
		accessTokenclaims["role"] = "admin"
	} else {
		accessTokenclaims["role"] = "user"
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenclaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenclaims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"type":  "refresh",
		"exp":   jwt.TimeFunc().Add(RefreshTokenDuration).Unix(),
	}

	if user.IsAdmin() {
		refreshTokenclaims["role"] = "admin"
	} else {
		refreshTokenclaims["role"] = "user"
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenclaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, err
	}

	return &domain.Token{
		Access:  accessTokenString,
		Refresh: refreshTokenString,
	}, nil
}

// ValidateToken implements JWTService.
func (j *JWTServiceImpl) ValidateToken(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// GenerateAccess implements JWTService.
func (j *JWTServiceImpl) GenerateAccess(claims jwt.MapClaims) (string, error) {
	accessTokenclaims := jwt.MapClaims{
		"id":    claims["id"],
		"email": claims["email"],
		"role":  claims["role"],
		"type":  "access",
		"exp":   jwt.TimeFunc().Add(AccessTokenDuration).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenclaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}
