/*
================================================================================
GATEWAY SERVER - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements the main HTTP server for the API gateway.

Go Concepts Used:
- net/http: HTTP server implementation
- JSON: API request/response handling
- Goroutines: Concurrent request handling
- Middleware: Authentication and logging
================================================================================
*/

package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"knowledge-exchange/analytics"
	"knowledge-exchange/auth"
	"knowledge-exchange/library"
	"knowledge-exchange/models"
	"knowledge-exchange/storage"
	"knowledge-exchange/utils"
)

// ============================================================================
// SERVER STRUCT
// ============================================================================

// Server represents the API gateway server
type Server struct {
	// HTTP server
	httpServer *http.Server

	// Authentication services
	authService *auth.Service
	userStore   *storage.UserStore

	// Services
	peerRegistry      *models.PeerRegistry
	fileIndex         *models.FileIndex
	indexer           *library.Indexer
	transferManager   *library.TransferManager
	integrityService  *library.IntegrityService
	reputationService *analytics.ReputationService
	ratingService     *analytics.RatingService
	throttlingManager *analytics.ThrottlingManager

	// Router
	router *Router

	// Discovery service
	discovery *Discovery

	// Server state
	isRunning bool
	mutex     sync.RWMutex

	// Configuration
	config *utils.Config
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewServer creates a new Gateway server
func NewServer(config *utils.Config) *Server {
	// Initialize authentication services
	authService := auth.NewService()
	userStore := storage.NewUserStore()

	// Initialize core data structures
	peerRegistry := models.NewPeerRegistry()
	fileIndex := models.NewFileIndex()

	// Initialize services
	indexer := library.NewIndexer(config.SharedFilesDir)
	transferManager := library.NewTransferManager(indexer)
	integrityService := library.NewIntegrityService()
	reputationService := analytics.NewReputationService(peerRegistry)
	ratingService := analytics.NewRatingService(reputationService)
	throttlingManager := analytics.NewThrottlingManager()
	discovery := NewDiscovery(peerRegistry)

	server := &Server{
		authService:       authService,
		userStore:         userStore,
		peerRegistry:      peerRegistry,
		fileIndex:         fileIndex,
		indexer:           indexer,
		transferManager:   transferManager,
		integrityService:  integrityService,
		reputationService: reputationService,
		ratingService:     ratingService,
		throttlingManager: throttlingManager,
		discovery:         discovery,
		isRunning:         false,
		config:            config,
	}

	// Create router with server reference
	server.router = NewRouter(server)

	return server
}

// ============================================================================
// SERVER LIFECYCLE
// ============================================================================

// Start starts the HTTP server
func (s *Server) Start() error {
	s.mutex.Lock()
	if s.isRunning {
		s.mutex.Unlock()
		return fmt.Errorf("server already running")
	}
	s.isRunning = true
	s.mutex.Unlock()

	// Start services
	s.reputationService.Start()
	s.ratingService.Start()
	s.discovery.Start()

	// Start file watcher
	s.indexer.StartWatcher(s.config.PeerID, 30*time.Second)

	// Configure HTTP server
	addr := fmt.Sprintf(":%d", s.config.APIPort)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router.GetHandler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server in goroutine
	go func() {
		log.Printf("Starting server on %s", addr)
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return nil
}

// Stop stops the HTTP server gracefully
func (s *Server) Stop() error {
	s.mutex.Lock()
	if !s.isRunning {
		s.mutex.Unlock()
		return nil
	}
	s.isRunning = false
	s.mutex.Unlock()

	// Stop services
	s.reputationService.Stop()
	s.ratingService.Stop()
	s.discovery.Stop()
	s.indexer.StopWatcher()
	s.throttlingManager.StopAll()

	// Shutdown HTTP server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// IsRunning returns the server running state
func (s *Server) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isRunning
}

// ============================================================================
// API RESPONSE TYPES
// ============================================================================

// APIResponse is the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PeerInfo contains public peer information
type PeerInfo struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Reputation float64 `json:"reputation"`
	IsOnline   bool    `json:"is_online"`
	Uploads    int     `json:"uploads"`
	Downloads  int     `json:"downloads"`
}

// FileInfo contains public file information
type FileInfo struct {
	CID        string    `json:"cid"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	Type       string    `json:"type"`
	Subject    string    `json:"subject"`
	OwnerID    string    `json:"owner_id"`
	Downloads  int       `json:"downloads"`
	Rating     float64   `json:"rating"`
	Available  bool      `json:"available"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// ============================================================================
// API HANDLERS
// ============================================================================

// HandleStatus returns server status
func (s *Server) HandleStatus(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status":     "running",
		"version":    utils.AppVersion,
		"peer_count": s.peerRegistry.Count(),
		"file_count": s.fileIndex.Count(),
		"uptime":     "active",
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// HandleRegister handles peer registration
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string `json:"name"`
		IPAddress string `json:"ip_address"`
		Port      int    `json:"port"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create new peer
	peerID := utils.GeneratePeerID(req.Name, req.IPAddress, req.Port)
	student := models.NewStudent(peerID, req.Name, req.IPAddress, req.Port)

	// Register
	s.peerRegistry.Register(student)

	s.sendJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Peer registered successfully",
		Data: PeerInfo{
			ID:         student.ID,
			Name:       student.Name,
			Reputation: student.ReputationScore,
			IsOnline:   student.IsOnline,
		},
	})
}

// HandleGetPeers returns list of online peers
func (s *Server) HandleGetPeers(w http.ResponseWriter, r *http.Request) {
	peers := s.peerRegistry.GetOnlinePeers()

	peerList := make([]PeerInfo, len(peers))
	for i, p := range peers {
		peerList[i] = PeerInfo{
			ID:         p.ID,
			Name:       p.Name,
			Reputation: p.ReputationScore,
			IsOnline:   p.IsOnline,
			Uploads:    p.TotalUploads,
			Downloads:  p.TotalDownloads,
		}
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    peerList,
	})
}

// HandleSearch handles file search requests
func (s *Server) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		s.sendError(w, http.StatusBadRequest, "Search query required")
		return
	}

	files := s.indexer.Search(query)

	fileList := make([]FileInfo, len(files))
	for i, f := range files {
		fileList[i] = FileInfo{
			CID:        f.CID,
			Name:       f.FileName,
			Size:       f.Size,
			Type:       f.FileType,
			Subject:    f.Subject,
			OwnerID:    f.OwnerID,
			Downloads:  f.DownloadCount,
			Rating:     f.AverageRating,
			Available:  f.IsAvailable,
			UploadedAt: f.UploadTime,
		}
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    fileList,
	})
}

// HandleGetFiles returns all available files
func (s *Server) HandleGetFiles(w http.ResponseWriter, r *http.Request) {
	files := s.indexer.GetAllFiles()

	fileList := make([]FileInfo, len(files))
	for i, f := range files {
		fileList[i] = FileInfo{
			CID:        f.CID,
			Name:       f.FileName,
			Size:       f.Size,
			Type:       f.FileType,
			Subject:    f.Subject,
			OwnerID:    f.OwnerID,
			Downloads:  f.DownloadCount,
			Rating:     f.AverageRating,
			Available:  f.IsAvailable,
			UploadedAt: f.UploadTime,
		}
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    fileList,
	})
}

// HandleGetReputation returns reputation for a peer
func (s *Server) HandleGetReputation(w http.ResponseWriter, r *http.Request) {
	peerID := r.URL.Query().Get("peer_id")
	if peerID == "" {
		s.sendError(w, http.StatusBadRequest, "Peer ID required")
		return
	}

	reputation, err := s.reputationService.GetReputation(peerID)
	if err != nil {
		s.sendError(w, http.StatusNotFound, err.Error())
		return
	}

	canDownload, reason := s.reputationService.CanDownload(peerID)

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"peer_id":      peerID,
			"reputation":   reputation,
			"can_download": canDownload,
			"reason":       reason,
		},
	})
}

// HandleRateFile handles file rating submission
func (s *Server) HandleRateFile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RaterID string  `json:"rater_id"`
		FileCID string  `json:"file_cid"`
		Score   float64 `json:"score"`
		Comment string  `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	rating, err := s.ratingService.RateFile(req.RaterID, req.FileCID, req.Score, req.Comment)
	if err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.sendJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Rating submitted successfully",
		Data:    rating,
	})
}

// HandleGetStats returns system statistics
func (s *Server) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"peers":      s.peerRegistry.Count(),
		"files":      s.fileIndex.Count(),
		"indexer":    s.indexer.GetStats(),
		"transfers":  s.transferManager.GetStats(),
		"reputation": s.reputationService.GetStats(),
		"ratings":    s.ratingService.GetGlobalStats(),
		"throttling": s.throttlingManager.GetStats(),
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// sendJSON sends a JSON response
func (s *Server) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// sendError sends an error response
func (s *Server) sendError(w http.ResponseWriter, status int, message string) {
	s.sendJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}

// ============================================================================
// GETTERS FOR SERVICES
// ============================================================================

func (s *Server) GetPeerRegistry() *models.PeerRegistry              { return s.peerRegistry }
func (s *Server) GetFileIndex() *models.FileIndex                    { return s.fileIndex }
func (s *Server) GetIndexer() *library.Indexer                       { return s.indexer }
func (s *Server) GetTransferManager() *library.TransferManager       { return s.transferManager }
func (s *Server) GetReputationService() *analytics.ReputationService { return s.reputationService }
func (s *Server) GetRatingService() *analytics.RatingService         { return s.ratingService }
func (s *Server) GetThrottlingManager() *analytics.ThrottlingManager { return s.throttlingManager }
func (s *Server) GetDiscovery() *Discovery                           { return s.discovery }
