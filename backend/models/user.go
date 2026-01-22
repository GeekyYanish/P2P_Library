/*
================================================================================
USER MODEL - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file defines the User model for authentication and authorization.

Go Concepts Used:
- Structs: Data models
- Methods: Business logic on models
- Time: Timestamp handling
- Validation: Input validation
================================================================================
*/

package models

import (
	"errors"
	"regexp"
	"time"
)

// ============================================================================
// CONSTANTS
// ============================================================================

// User roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// ============================================================================
// USER MODEL
// ============================================================================

// User represents a registered user in the system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never send password hash to client
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    time.Time `json:"last_login,omitempty"`
	IsActive     bool      `json:"is_active"`

	// P2P Network fields (from original Peer model)
	PeerID         string  `json:"peer_id,omitempty"`
	Reputation     float64 `json:"reputation"`
	TotalUploads   int     `json:"total_uploads"`
	TotalDownloads int     `json:"total_downloads"`
}

// ============================================================================
// VALIDATION METHODS
// ============================================================================

// Validate validates the user data
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	if u.Username == "" {
		return errors.New("username is required")
	}

	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if u.Role != RoleUser && u.Role != RoleAdmin {
		return errors.New("invalid role")
	}

	return nil
}

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	// Simple email regex pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
}

// CanDownload checks if user's reputation allows downloads
func (u *User) CanDownload() bool {
	return u.Reputation >= 3.0
}

// PublicUser returns a user object safe for public display (no sensitive info)
type PublicUser struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	Reputation     float64   `json:"reputation"`
	TotalUploads   int       `json:"total_uploads"`
	TotalDownloads int       `json:"total_downloads"`
}

// ToPublic converts a User to PublicUser
func (u *User) ToPublic() PublicUser {
	return PublicUser{
		ID:             u.ID,
		Email:          u.Email,
		Username:       u.Username,
		Role:           u.Role,
		CreatedAt:      u.CreatedAt,
		Reputation:     u.Reputation,
		TotalUploads:   u.TotalUploads,
		TotalDownloads: u.TotalDownloads,
	}
}
