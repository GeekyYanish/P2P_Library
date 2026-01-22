/*
================================================================================
AUTHENTICATION SERVICE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements JWT-based authentication and authorization.

Go Concepts Used:
- JWT tokens: Secure authentication
- bcrypt: Password hashing
- Middleware: Request authentication
- Error handling
================================================================================
*/

package auth

import (
	"errors"
	"time"

	"knowledge-exchange/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// JWT secret key (in production, this should be in environment variables)
	jwtSecret = "your-secret-key-change-this-in-production"

	// Token expiration times
	tokenExpiration   = 24 * time.Hour     // 24 hours
	refreshExpiration = 7 * 24 * time.Hour // 7 days
)

// ============================================================================
// JWT CLAIMS
// ============================================================================

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// ============================================================================
// AUTHENTICATION SERVICE
// ============================================================================

// Service provides authentication functionality
type Service struct {
	secret []byte
}

// NewService creates a new authentication service
func NewService() *Service {
	return &Service{
		secret: []byte(jwtSecret),
	}
}

// ============================================================================
// PASSWORD METHODS
// ============================================================================

// HashPassword hashes a plain text password
func (s *Service) HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password must be at least 6 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func (s *Service) VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ============================================================================
// TOKEN METHODS
// ============================================================================

// GenerateToken generates a JWT token for a user
func (s *Service) GenerateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractUserID extracts the user ID from a token
func (s *Service) ExtractUserID(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}

// ============================================================================
// AUTHORIZATION HELPERS
// ============================================================================

// IsAdmin checks if the user (from claims) is an admin
func IsAdmin(claims *Claims) bool {
	return claims.Role == models.RoleAdmin
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Add more password strength requirements as needed
	// For now, just checking minimum length

	return nil
}
