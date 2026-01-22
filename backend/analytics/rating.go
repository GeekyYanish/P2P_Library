/*
================================================================================
RATING SERVICE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements the rating aggregation system for files and peers.

Go Concepts Used:
- Maps: Storing ratings by target
- Slices: Dynamic rating lists
- Goroutines: Concurrent rating processing
- Channels: Async rating updates
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
	MinRatingValue = 1.0
	MaxRatingValue = 5.0
)

// ============================================================================
// RATING SERVICE STRUCT
// ============================================================================

// RatingService manages the rating system for files and peers
type RatingService struct {
	// ratingStore holds all ratings
	ratingStore *models.RatingStore

	// reputationService for updating reputation on ratings
	reputationService *ReputationService

	// ratingChan for async rating submissions
	ratingChan chan *models.Rating

	// mutex for thread-safe operations
	mutex sync.RWMutex

	// isRunning indicates if the service is active
	isRunning bool

	// stopChan signals the service to stop
	stopChan chan struct{}

	// Aggregated stats
	totalFileRatings  int
	totalPeerRatings  int
	averageFileRating float64
	averagePeerRating float64
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewRatingService creates a new RatingService
func NewRatingService(reputationService *ReputationService) *RatingService {
	return &RatingService{
		ratingStore:       models.NewRatingStore(),
		reputationService: reputationService,
		ratingChan:        make(chan *models.Rating, 100),
		isRunning:         false,
		stopChan:          make(chan struct{}),
	}
}

// ============================================================================
// SERVICE LIFECYCLE
// ============================================================================

// Start begins the rating processing goroutine
func (rs *RatingService) Start() {
	if rs.isRunning {
		return
	}

	rs.isRunning = true

	// Start rating processor goroutine
	go rs.processRatings()
}

// Stop stops the rating service
func (rs *RatingService) Stop() {
	if rs.isRunning {
		rs.isRunning = false
		close(rs.stopChan)
		close(rs.ratingChan)
	}
}

// processRatings handles rating submissions from the channel
func (rs *RatingService) processRatings() {
	for {
		select {
		case rating, ok := <-rs.ratingChan:
			if !ok {
				return
			}
			rs.applyRating(rating)
		case <-rs.stopChan:
			return
		}
	}
}

// ============================================================================
// RATING SUBMISSION
// ============================================================================

// RateFile submits a rating for a file
// Parameters:
//   - raterID: ID of the student giving the rating
//   - fileCID: Content ID of the file being rated
//   - score: Rating score (1-5)
//   - comment: Optional comment
//
// Returns:
//   - *models.Rating: The created rating
//   - error: Error if rating fails
func (rs *RatingService) RateFile(raterID, fileCID string, score float64, comment string) (*models.Rating, error) {
	// Validate score
	if score < MinRatingValue || score > MaxRatingValue {
		return nil, fmt.Errorf("score must be between %.0f and %.0f", MinRatingValue, MaxRatingValue)
	}

	// Check if already rated
	if rs.ratingStore.HasRated(raterID, fileCID) {
		return nil, fmt.Errorf("user has already rated this file")
	}

	// Create rating
	rating := models.NewRating(
		generateRatingID(raterID, fileCID),
		raterID,
		"file",
		fileCID,
		score,
		comment,
	)

	// Submit for async processing
	rs.ratingChan <- rating

	return rating, nil
}

// RatePeer submits a rating for another peer
func (rs *RatingService) RatePeer(raterID, targetPeerID string, score float64, comment string) (*models.Rating, error) {
	// Validate score
	if score < MinRatingValue || score > MaxRatingValue {
		return nil, fmt.Errorf("score must be between %.0f and %.0f", MinRatingValue, MaxRatingValue)
	}

	// Prevent self-rating
	if raterID == targetPeerID {
		return nil, fmt.Errorf("cannot rate yourself")
	}

	// Check if already rated
	if rs.ratingStore.HasRated(raterID, targetPeerID) {
		return nil, fmt.Errorf("user has already rated this peer")
	}

	// Create rating
	rating := models.NewRating(
		generateRatingID(raterID, targetPeerID),
		raterID,
		"student",
		targetPeerID,
		score,
		comment,
	)

	// Submit for async processing
	rs.ratingChan <- rating

	return rating, nil
}

// ============================================================================
// RATING APPLICATION
// ============================================================================

// applyRating processes a rating submission
func (rs *RatingService) applyRating(rating *models.Rating) error {
	// Add to store
	if err := rs.ratingStore.Add(rating); err != nil {
		return err
	}

	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// Update statistics
	if rating.TargetType == "file" {
		rs.totalFileRatings++
		rs.updateAverageFileRating(rating.Score)
	} else {
		rs.totalPeerRatings++
		rs.updateAveragePeerRating(rating.Score)

		// Update reputation for peer ratings
		if rs.reputationService != nil {
			rs.reputationService.RecordRating(rating.TargetID, rating.Score)
		}
	}

	return nil
}

// updateAverageFileRating recalculates average file rating
func (rs *RatingService) updateAverageFileRating(newScore float64) {
	// Incremental average calculation
	oldTotal := rs.averageFileRating * float64(rs.totalFileRatings-1)
	rs.averageFileRating = (oldTotal + newScore) / float64(rs.totalFileRatings)
}

// updateAveragePeerRating recalculates average peer rating
func (rs *RatingService) updateAveragePeerRating(newScore float64) {
	oldTotal := rs.averagePeerRating * float64(rs.totalPeerRatings-1)
	rs.averagePeerRating = (oldTotal + newScore) / float64(rs.totalPeerRatings)
}

// ============================================================================
// QUERY METHODS
// ============================================================================

// GetFileRatings returns all ratings for a file
func (rs *RatingService) GetFileRatings(fileCID string) []*models.Rating {
	return rs.ratingStore.GetByTarget(fileCID)
}

// GetPeerRatings returns all ratings for a peer
func (rs *RatingService) GetPeerRatings(peerID string) []*models.Rating {
	return rs.ratingStore.GetByTarget(peerID)
}

// GetFileStats returns rating statistics for a file
func (rs *RatingService) GetFileStats(fileCID string) models.RatingStats {
	return rs.ratingStore.GetStats(fileCID)
}

// GetPeerStats returns rating statistics for a peer
func (rs *RatingService) GetPeerStats(peerID string) models.RatingStats {
	return rs.ratingStore.GetStats(peerID)
}

// GetRatingsByRater returns all ratings given by a specific user
func (rs *RatingService) GetRatingsByRater(raterID string) []*models.Rating {
	return rs.ratingStore.GetByRater(raterID)
}

// HasUserRated checks if a user has rated a specific target
func (rs *RatingService) HasUserRated(raterID, targetID string) bool {
	return rs.ratingStore.HasRated(raterID, targetID)
}

// ============================================================================
// AGGREGATION METHODS
// ============================================================================

// AggregatedRating holds aggregated rating information
type AggregatedRating struct {
	TargetID     string    `json:"target_id"`
	TargetType   string    `json:"target_type"`
	TotalRatings int       `json:"total_ratings"`
	AverageScore float64   `json:"average_score"`
	Distribution [5]int    `json:"distribution"` // Count for each star (1-5)
	LastRatedAt  time.Time `json:"last_rated_at"`
}

// GetAggregatedRating returns aggregated rating info for a target
func (rs *RatingService) GetAggregatedRating(targetID, targetType string) *AggregatedRating {
	stats := rs.ratingStore.GetStats(targetID)

	return &AggregatedRating{
		TargetID:     targetID,
		TargetType:   targetType,
		TotalRatings: stats.TotalRatings,
		AverageScore: stats.AverageScore,
		Distribution: stats.RatingCounts,
		LastRatedAt:  stats.LastRatingAt,
	}
}

// GetTopRatedFiles returns the highest rated files
func (rs *RatingService) GetTopRatedFiles(fileIndex *models.FileIndex, limit int) []*models.AcademicFile {
	files := fileIndex.GetAllFiles()

	// Sort by average rating (bubble sort for simplicity)
	for i := 0; i < len(files)-1; i++ {
		for j := 0; j < len(files)-i-1; j++ {
			if files[j].AverageRating < files[j+1].AverageRating {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}

	if limit > len(files) {
		limit = len(files)
	}

	return files[:limit]
}

// ============================================================================
// STATISTICS
// ============================================================================

// GetGlobalStats returns global rating statistics
func (rs *RatingService) GetGlobalStats() map[string]interface{} {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return map[string]interface{}{
		"total_file_ratings":  rs.totalFileRatings,
		"total_peer_ratings":  rs.totalPeerRatings,
		"average_file_rating": rs.averageFileRating,
		"average_peer_rating": rs.averagePeerRating,
		"total_ratings":       rs.ratingStore.Count(),
		"is_running":          rs.isRunning,
	}
}

// ExportRatings exports all ratings as JSON
func (rs *RatingService) ExportRatings(targetID string) ([]byte, error) {
	ratings := rs.ratingStore.GetByTarget(targetID)
	return json.Marshal(ratings)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// generateRatingID creates a unique ID for a rating
func generateRatingID(raterID, targetID string) string {
	return fmt.Sprintf("rating-%s-%s-%d", raterID[:8], targetID[:8], time.Now().UnixNano())
}
