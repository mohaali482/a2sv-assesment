package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohaali482/a2sv-assesment/infrastructure"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	AdminMiddleware() gin.HandlerFunc
}

type MiddlewareImpl struct {
	jwtService infrastructure.JWTService
}

func NewMiddlewareImpl(jwtService infrastructure.JWTService) Middleware {
	return &MiddlewareImpl{
		jwtService: jwtService,
	}
}

func (m *MiddlewareImpl) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "auth header required"})
			c.Abort()
			return
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unautorized access"})
			c.Abort()
			return
		}
		tokenClaim, err := m.jwtService.ValidateToken(authParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unautorized access"})
			c.Abort()
			return
		}

		c.Set("user_id", tokenClaim["id"])
		c.Set("email", tokenClaim["email"])
		c.Set("role", tokenClaim["role"])
		c.Next()
	}
}

func (m *MiddlewareImpl) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unautorized access"})
			c.Abort()
			return
		}
		c.Next()
	}
}
