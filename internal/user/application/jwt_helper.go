package application

import (
	"time"

	"github.com/ddd-micro/internal/user/domain"
	"github.com/ddd-micro/pkg/security"
)

var (
	ErrInvalidToken = security.ErrInvalidToken
	ErrExpiredToken = security.ErrExpiredToken
)

// JWTClaims wraps security.JWTClaims with domain.Role
type JWTClaims struct {
	UserID uint        `json:"user_id"`
	Email  string      `json:"email"`
	Role   domain.Role `json:"role"`
}

// JWTHelper wraps security.JWTHelper for application use
type JWTHelper struct {
	helper *security.JWTHelper
}

// NewJWTHelper creates a new JWT helper
func NewJWTHelper(secretKey string, tokenDuration time.Duration) *JWTHelper {
	return &JWTHelper{
		helper: security.NewJWTHelper(secretKey, tokenDuration),
	}
}

// GenerateToken generates a new JWT token for a user
func (j *JWTHelper) GenerateToken(userID uint, email string, role domain.Role) (string, error) {
	return j.helper.GenerateToken(userID, email, role.String())
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTHelper) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims, err := j.helper.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &JWTClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   domain.Role(claims.Role),
	}, nil
}

// RefreshToken generates a new token with extended expiration
func (j *JWTHelper) RefreshToken(tokenString string) (string, error) {
	return j.helper.RefreshToken(tokenString)
}
