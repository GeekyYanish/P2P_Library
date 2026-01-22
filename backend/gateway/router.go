/*
================================================================================
ROUTER - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements API routing and middleware.

Go Concepts Used:
- http.Handler: Standard HTTP handler interface
- Middleware: Request preprocessing
- Mux: HTTP request routing
- Closures: Handler generation
================================================================================
*/

package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"knowledge-exchange/models"
)

// ============================================================================
// ROUTER STRUCT
// ============================================================================

// Router handles HTTP request routing
type Router struct {
	mux    *http.ServeMux
	server *Server
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewRouter creates a new Router
func NewRouter(server *Server) *Router {
	r := &Router{
		mux:    http.NewServeMux(),
		server: server,
	}

	r.setupRoutes()
	return r
}

// ============================================================================
// ROUTE SETUP
// ============================================================================

// setupRoutes configures all API routes
func (r *Router) setupRoutes() {
	// Health and status
	r.handle("GET", "/api/health", r.healthHandler())
	r.handle("GET", "/api/status", r.server.HandleStatus)

	// Authentication (public routes)
	r.handle("POST", "/api/auth/register", r.registerHandler())
	r.handle("POST", "/api/auth/login", r.loginHandler())
	r.handle("POST", "/api/auth/logout", r.logoutHandler())
	r.handle("GET", "/api/auth/me", r.meHandler())

	// Peer management
	r.handle("POST", "/api/peers/register", r.server.HandleRegister)
	r.handle("GET", "/api/peers", r.server.HandleGetPeers)
	r.handle("GET", "/api/peers/online", r.onlinePeersHandler())

	// File operations
	r.handle("GET", "/api/files", r.server.HandleGetFiles)
	r.handle("GET", "/api/files/search", r.server.HandleSearch)
	r.handle("POST", "/api/files/upload", r.uploadHandler())
	r.handle("GET", "/api/files/download", r.downloadHandler())

	// Reputation
	r.handle("GET", "/api/reputation", r.server.HandleGetReputation)
	r.handle("GET", "/api/reputation/history", r.reputationHistoryHandler())
	r.handle("GET", "/api/reputation/top", r.topContributorsHandler())

	// Ratings
	r.handle("POST", "/api/ratings/file", r.server.HandleRateFile)
	r.handle("POST", "/api/ratings/peer", r.ratePeerHandler())
	r.handle("GET", "/api/ratings", r.getRatingsHandler())

	// Statistics
	r.handle("GET", "/api/stats", r.server.HandleGetStats)

	// Static files (for frontend)
	r.mux.Handle("/", http.FileServer(http.Dir("../frontend")))
}

// handle registers a handler with method checking
func (r *Router) handle(method, pattern string, handler http.HandlerFunc) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		// Apply middleware chain
		finalHandler := r.applyMiddleware(handler)

		// Check method
		if req.Method != method && method != "" {
			// Allow OPTIONS for CORS
			if req.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.WriteHeader(http.StatusOK)
				return
			}
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		finalHandler.ServeHTTP(w, req)
	})
}

// GetHandler returns the http.Handler for the router
func (r *Router) GetHandler() http.Handler {
	return r.mux
}

// ============================================================================
// MIDDLEWARE
// ============================================================================

// applyMiddleware wraps a handler with all middleware
func (r *Router) applyMiddleware(handler http.HandlerFunc) http.Handler {
	// Apply middleware in reverse order (first applied runs first)
	h := http.Handler(handler)
	h = r.corsMiddleware(h)
	h = r.loggingMiddleware(h)
	h = r.recoveryMiddleware(h)
	return h
}

// loggingMiddleware logs all requests
func (r *Router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		// Create response writer wrapper to capture status
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Process request
		next.ServeHTTP(wrapped, req)

		// Log request
		log.Printf(
			"%s %s %d %v",
			req.Method,
			req.URL.Path,
			wrapped.status,
			time.Since(start),
		)
	})
}

// corsMiddleware adds CORS headers
func (r *Router) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, req)
	})
}

// recoveryMiddleware recovers from panics
func (r *Router) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, req)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// ============================================================================
// ADDITIONAL HANDLERS
// ============================================================================

// healthHandler returns a health check endpoint
func (r *Router) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	}
}

// onlinePeersHandler returns online peers
func (r *Router) onlinePeersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		peers := r.server.GetDiscovery().GetOnlinePeers()
		peerInfos := make([]PeerInfo, len(peers))

		for i, p := range peers {
			peerInfos[i] = PeerInfo{
				ID:         p.ID,
				Name:       p.Name,
				Reputation: p.ReputationScore,
				IsOnline:   p.IsOnline,
			}
		}

		r.server.sendJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    peerInfos,
		})
	}
}

// uploadHandler handles file upload
func (r *Router) uploadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Parse multipart form
		err := req.ParseMultipartForm(100 << 20) // 100 MB max
		if err != nil {
			r.server.sendError(w, http.StatusBadRequest, "Failed to parse form")
			return
		}

		// Get uploader ID
		ownerID := req.FormValue("owner_id")
		if ownerID == "" {
			r.server.sendError(w, http.StatusBadRequest, "Owner ID required")
			return
		}

		// Get file
		file, header, err := req.FormFile("file")
		if err != nil {
			r.server.sendError(w, http.StatusBadRequest, "File required")
			return
		}
		defer file.Close()

		// Read file content
		content := make([]byte, header.Size)
		_, err = file.Read(content)
		if err != nil {
			r.server.sendError(w, http.StatusInternalServerError, "Failed to read file")
			return
		}

		// Get file extension
		parts := strings.Split(header.Filename, ".")
		ext := ""
		if len(parts) > 1 {
			ext = "." + parts[len(parts)-1]
		}

		// Create academic file
		academicFile := models.NewAcademicFile(header.Filename, ownerID, header.Size, ext, content)

		// Add to index
		r.server.GetFileIndex().Add(academicFile)

		// Record upload for reputation
		r.server.GetReputationService().RecordUpload(ownerID)

		r.server.sendJSON(w, http.StatusCreated, APIResponse{
			Success: true,
			Message: "File uploaded successfully",
			Data: map[string]interface{}{
				"cid":       academicFile.CID,
				"file_name": academicFile.FileName,
				"size":      academicFile.Size,
			},
		})
	}
}

// downloadHandler handles file download
func (r *Router) downloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		cid := req.URL.Query().Get("cid")
		requesterID := req.URL.Query().Get("requester_id")

		if cid == "" || requesterID == "" {
			r.server.sendError(w, http.StatusBadRequest, "CID and requester_id required")
			return
		}

		// Check reputation
		canDownload, reason := r.server.GetReputationService().CanDownload(requesterID)
		if !canDownload {
			r.server.sendError(w, http.StatusForbidden, reason)
			return
		}

		// Get file
		file, exists := r.server.GetFileIndex().Get(cid)
		if !exists {
			r.server.sendError(w, http.StatusNotFound, "File not found")
			return
		}

		// Record download
		r.server.GetReputationService().RecordDownload(requesterID)
		file.RecordDownload()

		r.server.sendJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Message: "Download initiated",
			Data: map[string]interface{}{
				"cid":       file.CID,
				"file_name": file.FileName,
				"size":      file.Size,
				"owner":     file.OwnerID,
			},
		})
	}
}

// reputationHistoryHandler returns reputation history
func (r *Router) reputationHistoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		peerID := req.URL.Query().Get("peer_id")
		if peerID == "" {
			r.server.sendError(w, http.StatusBadRequest, "Peer ID required")
			return
		}

		history := r.server.GetReputationService().GetEventHistory(peerID)
		r.server.sendJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    history,
		})
	}
}

// topContributorsHandler returns top contributors
func (r *Router) topContributorsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		contributors := r.server.GetReputationService().GetTopContributors(10)

		peerInfos := make([]PeerInfo, len(contributors))
		for i, p := range contributors {
			peerInfos[i] = PeerInfo{
				ID:         p.ID,
				Name:       p.Name,
				Reputation: p.ReputationScore,
				Uploads:    p.TotalUploads,
				Downloads:  p.TotalDownloads,
			}
		}

		r.server.sendJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    peerInfos,
		})
	}
}

// ratePeerHandler handles peer rating
func (r *Router) ratePeerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			RaterID  string  `json:"rater_id"`
			TargetID string  `json:"target_id"`
			Score    float64 `json:"score"`
			Comment  string  `json:"comment"`
		}

		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			r.server.sendError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		rating, err := r.server.GetRatingService().RatePeer(
			body.RaterID, body.TargetID, body.Score, body.Comment,
		)
		if err != nil {
			r.server.sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		r.server.sendJSON(w, http.StatusCreated, APIResponse{
			Success: true,
			Message: "Peer rated successfully",
			Data:    rating,
		})
	}
}

// getRatingsHandler returns ratings for a target
func (r *Router) getRatingsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		targetID := req.URL.Query().Get("target_id")
		targetType := req.URL.Query().Get("type")

		if targetID == "" {
			r.server.sendError(w, http.StatusBadRequest, "Target ID required")
			return
		}

		var ratings interface{}
		if targetType == "file" {
			ratings = r.server.GetRatingService().GetFileRatings(targetID)
		} else {
			ratings = r.server.GetRatingService().GetPeerRatings(targetID)
		}

		r.server.sendJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    ratings,
		})
	}
}
