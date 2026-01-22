/*
================================================================================
AUTHENTICATION HANDLERS - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements HTTP handlers for authentication endpoints.

Go Concepts Used:
- HTTP handlers: Request/response handling
- JSON encoding/decoding
- Error handling
- Authentication
================================================================================
*/

package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"knowledge-exchange/auth"
	"knowledge-exchange/models"
)

// ============================================================================
// REQUEST/RESPONSE TYPES
// ============================================================================

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	Data    *AuthData `json:"data,omitempty"`
	Error   string    `json:"error,omitempty"`
}

// AuthData contains user data and token
type AuthData struct {
	User  models.PublicUser `json:"user"`
	Token string            `json:"token"`
}

// ============================================================================
// REGISTER HANDLER
// ============================================================================

// registerHandler handles user registration
func (r *Router) registerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var regReq RegisterRequest

		if err := json.NewDecoder(req.Body).Decode(&regReq); err != nil {
			sendJSON(w, http.StatusBadRequest, AuthResponse{
				Success: false,
				Error:   "Invalid request body",
			})
			return
		}

		// Validate password strength
		if err := auth.ValidatePasswordStrength(regReq.Password); err != nil {
			sendJSON(w, http.StatusBadRequest, AuthResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		// Hash password
		passwordHash, err := r.server.authService.HashPassword(regReq.Password)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, AuthResponse{
				Success: false,
				Error:   "Failed to process password",
			})
			return
		}

		// Create user
		user := &models.User{
			Email:        regReq.Email,
			Username:     regReq.Username,
			PasswordHash: passwordHash,
			Role:         models.RoleUser, // Default role
		}

		// Save user
		if err := r.server.userStore.Create(user); err != nil {
			sendJSON(w, http.StatusBadRequest, AuthResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		log.Printf("New user registered: %s (%s)", user.Username, user.Email)

		sendJSON(w, http.StatusCreated, AuthResponse{
			Success: true,
			Message: "Registration successful. Please login.",
		})
	}
}

// ============================================================================
// LOGIN HANDLER
// ============================================================================

// loginHandler handles user login
func (r *Router) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var loginReq LoginRequest

		if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
			sendJSON(w, http.StatusBadRequest, AuthResponse{
				Success: false,
				Error:   "Invalid request body",
			})
			return
		}

		// Find user by email
		user, err := r.server.userStore.GetByEmail(loginReq.Email)
		if err != nil {
			sendJSON(w, http.StatusUnauthorized, AuthResponse{
				Success: false,
				Error:   "Invalid email or password",
			})
			return
		}

		// Check if user is active
		if !user.IsActive {
			sendJSON(w, http.StatusUnauthorized, AuthResponse{
				Success: false,
				Error:   "Account is deactivated",
			})
			return
		}

		// Verify password
		if err := r.server.authService.VerifyPassword(user.PasswordHash, loginReq.Password); err != nil {
			sendJSON(w, http.StatusUnauthorized, AuthResponse{
				Success: false,
				Error:   "Invalid email or password",
			})
			return
		}

		// Generate JWT token
		token, err := r.server.authService.GenerateToken(user)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, AuthResponse{
				Success: false,
				Error:   "Failed to generate token",
			})
			return
		}

		// Update last login
		user.UpdateLastLogin()
		r.server.userStore.Update(user)

		log.Printf("User logged in: %s (%s)", user.Username, user.Email)

		sendJSON(w, http.StatusOK, AuthResponse{
			Success: true,
			Message: "Login successful",
			Data: &AuthData{
				User:  user.ToPublic(),
				Token: token,
			},
		})
	}
}

// ============================================================================
// GET CURRENT USER HANDLER
// ============================================================================

// meHandler returns the current user from the token
func (r *Router) meHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Extract token from header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			sendJSON(w, http.StatusUnauthorized, AuthResponse{
				Success: false,
				Error:   "Authorization header required",
			})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendJSON(w, http.StatusUnauthorized, AuthResponse{
				Success: false,
				Error:   "Invalid authorization header format",
			})
			return
		}

		token := parts[1]

		// Validate token and get claims
		claims, err := r.server.authService.ValidateToken(token)
		if err != nil {
			sendJSON(w, http.StatusUnauthorized, AuthResponse{
				Success: false,
				Error:   "Invalid or expired token",
			})
			return
		}

		// Get user from store
		user, err := r.server.userStore.GetByID(claims.UserID)
		if err != nil {
			sendJSON(w, http.StatusNotFound, AuthResponse{
				Success: false,
				Error:   "User not found",
			})
			return
		}

		sendJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    user.ToPublic(),
		})
	}
}

// ============================================================================
// LOGOUT HANDLER
// ============================================================================

// logoutHandler handles user logout (client-side token deletion)
func (r *Router) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// In JWT-based auth, logout is primarily client-side
		// The client should delete the token from storage
		// We could implement a token blacklist here if needed

		sendJSON(w, http.StatusOK, AuthResponse{
			Success: true,
			Message: "Logout successful",
		})
	}
}

// ============================================================================
// MIDDLEWARE
// ============================================================================

// authMiddleware validates JWT token and adds user info to request context
func (r *Router) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Extract token from header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			sendJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Authorization header required",
			})
			return
		}

		//Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Invalid authorization header format",
			})
			return
		}

		token := parts[1]

		// Validate token
		claims, err := r.server.authService.ValidateToken(token)
		if err != nil {
			sendJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Invalid or expired token",
			})
			return
		}

		// Store user ID in request context for handlers to use
		// For now, we'll validate the user exists
		_, err = r.server.userStore.GetByID(claims.UserID)
		if err != nil {
			sendJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "User not found",
			})
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, req)
	})
}

// adminMiddleware checks if user has admin role
func (r *Router) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Extract and validate token (similar to authMiddleware)
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			sendJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Authorization header required",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Invalid authorization header format",
			})
			return
		}

		token := parts[1]
		claims, err := r.server.authService.ValidateToken(token)
		if err != nil {
			sendJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Invalid or expired token",
			})
			return
		}

		// Check if user is admin
		if !auth.IsAdmin(claims) {
			sendJSON(w, http.StatusForbidden, map[string]interface{}{
				"success": false,
				"error":   "Admin access required",
			})
			return
		}

		next.ServeHTTP(w, req)
	})
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
