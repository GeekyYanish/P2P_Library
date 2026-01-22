/*
================================================================================
USER STORAGE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements in-memory storage for user data.

Go Concepts Used:
- Maps: In-memory data storage
- Sync.RWMutex: Thread-safe operations
- UUID: Unique identifiers
================================================================================
*/

package storage

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"knowledge-exchange/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// USER STORE
// ============================================================================

// UserStore manages user data in memory
type UserStore struct {
	users      map[string]*models.User // userID -> User
	emailIndex map[string]string       // email -> userID (for lookups)
	mu         sync.RWMutex
}

// NewUserStore creates a new user store
func NewUserStore() *UserStore {
	store := &UserStore{
		users:      make(map[string]*models.User),
		emailIndex: make(map[string]string),
	}

	// Create default admin user
	store.createDefaultAdmin()

	return store
}

// createDefaultAdmin creates a default admin user for testing
func (s *UserStore) createDefaultAdmin() {
	adminID := uuid.New().String()

	// Generate password hash for "admin123" at runtime
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Failed to create default admin user: %v", err)
		return
	}

	admin := &models.User{
		ID:           adminID,
		Email:        "admin@knowledge-exchange.com",
		Username:     "admin",
		PasswordHash: string(passwordHash),
		Role:         models.RoleAdmin,
		CreatedAt:    time.Now(),
		IsActive:     true,
		Reputation:   10.0,
	}

	s.users[adminID] = admin
	s.emailIndex[strings.ToLower(admin.Email)] = adminID
	log.Printf("✓ Default admin user created (admin@knowledge-exchange.com / admin123)")
}

// ============================================================================
// CRUD OPERATIONS
// ============================================================================

// Create creates a new user
func (s *UserStore) Create(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if email already exists
	emailLower := strings.ToLower(user.Email)
	if _, exists := s.emailIndex[emailLower]; exists {
		return errors.New("email already registered")
	}

	// Validate user
	if err := user.Validate(); err != nil {
		return err
	}

	// Generate ID if not set
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Set defaults
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.Reputation == 0 {
		user.Reputation = 5.0 // Default starting reputation
	}
	user.IsActive = true

	// Store user
	s.users[user.ID] = user
	s.emailIndex[emailLower] = user.ID

	return nil
}

// GetByID retrieves a user by ID
func (s *UserStore) GetByID(userID string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	emailLower := strings.ToLower(email)
	userID, exists := s.emailIndex[emailLower]
	if !exists {
		return nil, errors.New("user not found")
	}

	user := s.users[userID]
	return user, nil
}

// Update updates a user
func (s *UserStore) Update(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Validate user
	if err := user.Validate(); err != nil {
		return err
	}

	s.users[user.ID] = user

	return nil
}

// Delete deletes a user (soft delete by setting IsActive to false)
func (s *UserStore) Delete(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	user.IsActive = false

	return nil
}

// List returns all users
func (s *UserStore) List() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		if user.IsActive {
			users = append(users, user)
		}
	}

	return users
}

// Count returns the total number of active users
func (s *UserStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, user := range s.users {
		if user.IsActive {
			count++
		}
	}

	return count
}

// ============================================================================
// ADMIN OPERATIONS
// ============================================================================

// UpdateRole updates a user's role (admin only operation)
func (s *UserStore) UpdateRole(userID, newRole string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	if newRole != models.RoleUser && newRole != models.RoleAdmin {
		return errors.New("invalid role")
	}

	user.Role = newRole

	return nil
}

// UpdateReputation updates a user's reputation
func (s *UserStore) UpdateReputation(userID string, delta float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	user.Reputation += delta

	// Ensure reputation stays within bounds
	if user.Reputation < 0 {
		user.Reputation = 0
	}
	if user.Reputation > 10 {
		user.Reputation = 10
	}

	return nil
}
