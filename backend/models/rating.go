/*
================================================================================
RATING MODEL - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file defines the Rating struct for the peer rating system.

Go Concepts Used:
- Structs: For rating data structure
- Interfaces: Defining behavior contracts
- Slices: Storing rating history
- Maps: Aggregating ratings by file/student
================================================================================
*/

package models

import (
	"encoding/json"
	"sync"
	"time"
)

// ============================================================================
// CONSTANTS - Rating constraints
// ============================================================================

const (
	MinRating = 1.0 // Minimum rating value
	MaxRating = 5.0 // Maximum rating value
)

// ============================================================================
// STRUCT DEFINITIONS
// ============================================================================

// Rating represents a rating given by one student to a file or another student
type Rating struct {
	// ID is the unique identifier for this rating
	ID string `json:"id"`

	// RaterID is the student who gave the rating
	RaterID string `json:"rater_id"`

	// TargetType indicates what was rated ("file" or "student")
	TargetType string `json:"target_type"`

	// TargetID is the CID (for files) or student ID being rated
	TargetID string `json:"target_id"`

	// Score is the rating value (1-5)
	Score float64 `json:"score"`

	// Comment is an optional text review
	Comment string `json:"comment"`

	// Timestamp records when the rating was given
	Timestamp time.Time `json:"timestamp"`
}

// RatingStats holds aggregated rating statistics
type RatingStats struct {
	TotalRatings int       `json:"total_ratings"`
	AverageScore float64   `json:"average_score"`
	RatingCounts [5]int    `json:"rating_counts"` // Array of counts for ratings 1-5
	LastRatingAt time.Time `json:"last_rating_at"`
}

// ============================================================================
// INTERFACE DEFINITION
// ============================================================================

// Rateable is an interface for objects that can be rated
// Any type that implements these methods satisfies the Rateable interface
type Rateable interface {
	GetID() string
	AddRating(rating float64)
	GetAverageRating() float64
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewRating creates a new rating instance
func NewRating(id, raterID, targetType, targetID string, score float64, comment string) *Rating {
	// Validate score range
	if score < MinRating {
		score = MinRating
	}
	if score > MaxRating {
		score = MaxRating
	}

	return &Rating{
		ID:         id,
		RaterID:    raterID,
		TargetType: targetType,
		TargetID:   targetID,
		Score:      score,
		Comment:    comment,
		Timestamp:  time.Now(),
	}
}

// ============================================================================
// METHODS
// ============================================================================

// IsValid checks if the rating is valid
func (r *Rating) IsValid() (bool, string) {
	if r.RaterID == "" {
		return false, "Rater ID is required"
	}

	if r.TargetID == "" {
		return false, "Target ID is required"
	}

	if r.TargetType != "file" && r.TargetType != "student" {
		return false, "Target type must be 'file' or 'student'"
	}

	if r.Score < MinRating || r.Score > MaxRating {
		return false, "Score must be between 1 and 5"
	}

	// Prevent self-rating
	if r.RaterID == r.TargetID {
		return false, "Cannot rate yourself"
	}

	return true, ""
}

// ToJSON converts the rating to JSON bytes
func (r *Rating) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON populates the rating from JSON bytes
func (r *Rating) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// ============================================================================
// RATING STORE - Storage and aggregation for ratings
// ============================================================================

// RatingStore manages all ratings in the system
type RatingStore struct {
	// ratings stores all ratings by their ID
	ratings map[string]*Rating

	// byTarget groups ratings by their target ID for quick lookup
	byTarget map[string][]*Rating

	// mutex provides thread-safe access
	mutex sync.RWMutex
}

// NewRatingStore creates a new rating store
func NewRatingStore() *RatingStore {
	return &RatingStore{
		ratings:  make(map[string]*Rating),
		byTarget: make(map[string][]*Rating),
	}
}

// Add adds a new rating to the store
func (rs *RatingStore) Add(rating *Rating) error {
	// Validate the rating
	if valid, msg := rating.IsValid(); !valid {
		return &RatingError{Message: msg}
	}

	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// Store by ID
	rs.ratings[rating.ID] = rating

	// Add to target's rating list
	rs.byTarget[rating.TargetID] = append(rs.byTarget[rating.TargetID], rating)

	return nil
}

// GetByTarget returns all ratings for a specific target
func (rs *RatingStore) GetByTarget(targetID string) []*Rating {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return rs.byTarget[targetID]
}

// GetStats calculates aggregate statistics for a target
func (rs *RatingStore) GetStats(targetID string) RatingStats {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	ratings := rs.byTarget[targetID]
	stats := RatingStats{}

	if len(ratings) == 0 {
		return stats
	}

	var totalScore float64
	for _, r := range ratings {
		totalScore += r.Score
		stats.TotalRatings++

		// Count ratings by score (1-5 maps to index 0-4)
		index := int(r.Score) - 1
		if index >= 0 && index < 5 {
			stats.RatingCounts[index]++
		}

		// Track most recent rating
		if r.Timestamp.After(stats.LastRatingAt) {
			stats.LastRatingAt = r.Timestamp
		}
	}

	stats.AverageScore = totalScore / float64(stats.TotalRatings)
	return stats
}

// GetByRater returns all ratings given by a specific student
func (rs *RatingStore) GetByRater(raterID string) []*Rating {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	var raterRatings []*Rating
	for _, rating := range rs.ratings {
		if rating.RaterID == raterID {
			raterRatings = append(raterRatings, rating)
		}
	}
	return raterRatings
}

// HasRated checks if a rater has already rated a specific target
func (rs *RatingStore) HasRated(raterID, targetID string) bool {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	ratings := rs.byTarget[targetID]
	for _, r := range ratings {
		if r.RaterID == raterID {
			return true
		}
	}
	return false
}

// Count returns the total number of ratings
func (rs *RatingStore) Count() int {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return len(rs.ratings)
}

// ============================================================================
// ERROR TYPES
// ============================================================================

// RatingError represents a rating-related error
type RatingError struct {
	Message string
}

// Error implements the error interface
func (e *RatingError) Error() string {
	return e.Message
}
