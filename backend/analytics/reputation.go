/*
================================================================================
REPUTATION SERVICE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements the reputation calculation engine for fair access control.

Go Concepts Used:
- Goroutines: Concurrent reputation updates
- Channels: Asynchronous event processing
- Interfaces: Reputation source abstraction
- Mutex: Thread-safe reputation storage
================================================================================
*/

package analytics

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"knowledge-exchange/models"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// Reputation thresholds
	MinReputation     = 0.0
	MaxReputation     = 10.0
	DefaultReputation = 5.0
	DownloadThreshold = 3.0

	// Reputation change values
	UploadBonus      = 0.5
	DownloadPenalty  = 0.1
	GoodRatingBonus  = 0.3
	BadRatingPenalty = 0.2
	LeecherPenalty   = 0.5
	InactivityDecay  = 0.1
)

// ============================================================================
// REPUTATION EVENT TYPES
// ============================================================================

// ReputationEvent represents an event that affects reputation
type ReputationEvent struct {
	Type      string    `json:"type"`
	StudentID string    `json:"student_id"`
	Delta     float64   `json:"delta"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// Event types
const (
	EventUpload       = "UPLOAD"
	EventDownload     = "DOWNLOAD"
	EventRating       = "RATING"
	EventLeeching     = "LEECHING"
	EventInactivity   = "INACTIVITY"
	EventContribution = "CONTRIBUTION"
)

// ============================================================================
// REPUTATION SERVICE STRUCT
// ============================================================================

// ReputationService manages reputation calculations and enforcement
type ReputationService struct {
	// peerRegistry for accessing peer data
	peerRegistry *models.PeerRegistry

	// eventChan receives reputation events for processing
	eventChan chan ReputationEvent

	// eventHistory stores all reputation events
	eventHistory []ReputationEvent

	// mutex for thread-safe operations
	mutex sync.RWMutex

	// isRunning indicates if the service is active
	isRunning bool

	// stopChan signals the service to stop
	stopChan chan struct{}
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewReputationService creates a new ReputationService
func NewReputationService(peerRegistry *models.PeerRegistry) *ReputationService {
	return &ReputationService{
		peerRegistry: peerRegistry,
		eventChan:    make(chan ReputationEvent, 100),
		eventHistory: make([]ReputationEvent, 0),
		isRunning:    false,
		stopChan:     make(chan struct{}),
	}
}

// ============================================================================
// SERVICE LIFECYCLE
// ============================================================================

// Start begins the reputation processing goroutine
func (rs *ReputationService) Start() {
	if rs.isRunning {
		return
	}

	rs.isRunning = true

	// Start event processor goroutine
	go rs.processEvents()

	// Start decay checker goroutine
	go rs.checkInactivityDecay()
}

// Stop stops the reputation service
func (rs *ReputationService) Stop() {
	if rs.isRunning {
		rs.isRunning = false
		close(rs.stopChan)
		close(rs.eventChan)
	}
}

// processEvents handles reputation events from the channel
func (rs *ReputationService) processEvents() {
	for {
		select {
		case event, ok := <-rs.eventChan:
			if !ok {
				return
			}
			rs.applyEvent(event)
		case <-rs.stopChan:
			return
		}
	}
}

// checkInactivityDecay periodically applies decay to inactive peers
func (rs *ReputationService) checkInactivityDecay() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rs.applyInactivityDecay()
		case <-rs.stopChan:
			return
		}
	}
}

// ============================================================================
// REPUTATION CALCULATION
// ============================================================================

// CalculateReputation calculates current reputation for a student
// Based on upload/download ratio, ratings, and contributions
func (rs *ReputationService) CalculateReputation(student *models.Student) float64 {
	base := student.ReputationScore

	// Apply upload/download ratio factor
	if student.TotalDownloads > 0 {
		ratio := float64(student.TotalUploads) / float64(student.TotalDownloads)
		if ratio < 0.5 {
			// Downloading too much without uploading
			base -= 1.0
		} else if ratio > 1.5 {
			// Contributing more than taking
			base += 0.5
		}
	}

	// Clamp to valid range
	if base < MinReputation {
		base = MinReputation
	}
	if base > MaxReputation {
		base = MaxReputation
	}

	return base
}

// CanDownload checks if a student has sufficient reputation
func (rs *ReputationService) CanDownload(studentID string) (bool, string) {
	student, exists := rs.peerRegistry.Get(studentID)
	if !exists {
		return false, "Student not found"
	}

	reputation := rs.CalculateReputation(student)

	if reputation < DownloadThreshold {
		return false, fmt.Sprintf("Insufficient reputation (%.1f < %.1f)", reputation, DownloadThreshold)
	}

	return true, "Allowed"
}

// ============================================================================
// EVENT SUBMISSION
// ============================================================================

// RecordUpload records an upload event
func (rs *ReputationService) RecordUpload(studentID string) {
	event := ReputationEvent{
		Type:      EventUpload,
		StudentID: studentID,
		Delta:     UploadBonus,
		Reason:    "Shared a file",
		Timestamp: time.Now(),
	}
	rs.eventChan <- event
}

// RecordDownload records a download event
func (rs *ReputationService) RecordDownload(studentID string) {
	event := ReputationEvent{
		Type:      EventDownload,
		StudentID: studentID,
		Delta:     -DownloadPenalty,
		Reason:    "Downloaded a file",
		Timestamp: time.Now(),
	}
	rs.eventChan <- event
}

// RecordRating records a rating event
func (rs *ReputationService) RecordRating(studentID string, ratingScore float64) {
	var delta float64
	var reason string

	if ratingScore >= 4.0 {
		delta = GoodRatingBonus
		reason = "Received good rating"
	} else if ratingScore <= 2.0 {
		delta = -BadRatingPenalty
		reason = "Received bad rating"
	} else {
		return // Neutral rating, no effect
	}

	event := ReputationEvent{
		Type:      EventRating,
		StudentID: studentID,
		Delta:     delta,
		Reason:    reason,
		Timestamp: time.Now(),
	}
	rs.eventChan <- event
}

// RecordLeeching records a leeching penalty
func (rs *ReputationService) RecordLeeching(studentID string) {
	event := ReputationEvent{
		Type:      EventLeeching,
		StudentID: studentID,
		Delta:     -LeecherPenalty,
		Reason:    "Detected as leecher",
		Timestamp: time.Now(),
	}
	rs.eventChan <- event
}

// ============================================================================
// EVENT APPLICATION
// ============================================================================

// applyEvent applies a reputation event to a student
func (rs *ReputationService) applyEvent(event ReputationEvent) {
	student, exists := rs.peerRegistry.Get(event.StudentID)
	if !exists {
		return
	}

	// Apply the reputation change
	student.UpdateReputation(event.Delta)

	// Record in history
	rs.mutex.Lock()
	rs.eventHistory = append(rs.eventHistory, event)
	rs.mutex.Unlock()
}

// applyInactivityDecay applies reputation decay to inactive peers
func (rs *ReputationService) applyInactivityDecay() {
	peers := rs.peerRegistry.GetAllPeers()
	inactiveThreshold := 24 * time.Hour

	for _, peer := range peers {
		if time.Since(peer.LastSeen) > inactiveThreshold {
			event := ReputationEvent{
				Type:      EventInactivity,
				StudentID: peer.ID,
				Delta:     -InactivityDecay,
				Reason:    "Extended inactivity",
				Timestamp: time.Now(),
			}
			rs.applyEvent(event)
		}
	}
}

// ============================================================================
// QUERY METHODS
// ============================================================================

// GetReputation returns the current reputation for a student
func (rs *ReputationService) GetReputation(studentID string) (float64, error) {
	student, exists := rs.peerRegistry.Get(studentID)
	if !exists {
		return 0, fmt.Errorf("student not found: %s", studentID)
	}
	return rs.CalculateReputation(student), nil
}

// GetEventHistory returns reputation events for a student
func (rs *ReputationService) GetEventHistory(studentID string) []ReputationEvent {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	var events []ReputationEvent
	for _, event := range rs.eventHistory {
		if event.StudentID == studentID {
			events = append(events, event)
		}
	}
	return events
}

// GetTopContributors returns students with highest reputation
func (rs *ReputationService) GetTopContributors(limit int) []*models.Student {
	peers := rs.peerRegistry.GetAllPeers()

	// Sort by reputation (simple bubble sort for demonstration)
	for i := 0; i < len(peers)-1; i++ {
		for j := 0; j < len(peers)-i-1; j++ {
			if peers[j].ReputationScore < peers[j+1].ReputationScore {
				peers[j], peers[j+1] = peers[j+1], peers[j]
			}
		}
	}

	if limit > len(peers) {
		limit = len(peers)
	}

	return peers[:limit]
}

// GetLeechers returns all students marked as leechers
func (rs *ReputationService) GetLeechers() []*models.Student {
	peers := rs.peerRegistry.GetAllPeers()
	var leechers []*models.Student

	for _, peer := range peers {
		if peer.IsLeecher {
			leechers = append(leechers, peer)
		}
	}

	return leechers
}

// ============================================================================
// SERIALIZATION
// ============================================================================

// ExportHistory exports reputation history as JSON
func (rs *ReputationService) ExportHistory() ([]byte, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return json.Marshal(rs.eventHistory)
}

// GetStats returns reputation service statistics
func (rs *ReputationService) GetStats() map[string]interface{} {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return map[string]interface{}{
		"total_events":     len(rs.eventHistory),
		"is_running":       rs.isRunning,
		"top_contributors": len(rs.GetTopContributors(10)),
		"leecher_count":    len(rs.GetLeechers()),
	}
}
