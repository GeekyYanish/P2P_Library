/*
================================================================================
STUDENT MODEL - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file defines the Student struct representing a peer in the P2P network.

Go Concepts Used:
- Structs: Custom data types that group related fields together
- Pointers: Used to modify struct values across different functions
- Methods: Functions attached to structs using receiver syntax
- Type Inference: Using := for local variable initialization
- Constants: Defining fixed values like reputation thresholds
================================================================================
*/

package models

import (
	"encoding/json"
	"sync"
	"time"
)

// ============================================================================
// CONSTANTS - Define system thresholds
// ============================================================================

// MinReputation is the minimum reputation score required to download files
// Students below this threshold are considered "leechers" and face restrictions
const MinReputation float64 = 3.0

// MaxReputation is the maximum reputation score a student can achieve
const MaxReputation float64 = 10.0

// DefaultReputation is the starting reputation for new peers
const DefaultReputation float64 = 5.0

// ============================================================================
// STRUCT DEFINITION - The core Student data type
// ============================================================================

// Student represents a peer in the P2P academic network
// Each student has a unique ID and reputation score that determines their access level
type Student struct {
	// ID is the unique identifier for this peer (e.g., "node_1", "peer_abc123")
	ID string `json:"id"`

	// Name is the display name of the student
	Name string `json:"name"`

	// ReputationScore determines download privileges (range: 0.0 to 10.0)
	// Higher scores = more privileges, lower scores = throttling/restrictions
	ReputationScore float64 `json:"reputation_score"`

	// IsLeecher indicates if the student downloads more than they upload
	// Leechers face bandwidth restrictions to encourage fair sharing
	IsLeecher bool `json:"is_leecher"`

	// IsOnline indicates the current connection status
	IsOnline bool `json:"is_online"`

	// LastSeen records when the peer was last active
	LastSeen time.Time `json:"last_seen"`

	// TotalUploads tracks the number of files this peer has shared
	TotalUploads int `json:"total_uploads"`

	// TotalDownloads tracks the number of files this peer has downloaded
	TotalDownloads int `json:"total_downloads"`

	// IPAddress stores the peer's network address for P2P connections
	IPAddress string `json:"ip_address"`

	// Port is the port number this peer listens on
	Port int `json:"port"`
}

// ============================================================================
// CONSTRUCTOR - Create new Student instances
// ============================================================================

// NewStudent creates a new Student with default values
// This is a factory function pattern common in Go
// Parameters:
//   - id: Unique identifier for the student
//   - name: Display name
//   - ipAddress: Network address
//   - port: Listening port
//
// Returns:
//   - *Student: Pointer to the newly created student
func NewStudent(id, name, ipAddress string, port int) *Student {
	// Using := for type inference - Go automatically determines the type
	return &Student{
		ID:              id,
		Name:            name,
		ReputationScore: DefaultReputation, // Start with default reputation
		IsLeecher:       false,             // Not a leecher by default
		IsOnline:        true,              // Assume online when created
		LastSeen:        time.Now(),
		TotalUploads:    0,
		TotalDownloads:  0,
		IPAddress:       ipAddress,
		Port:            port,
	}
}

// ============================================================================
// METHODS - Functions attached to the Student struct
// ============================================================================

// CanDownload checks if the student has sufficient reputation to download files
// This uses the pointer receiver (*Student) to access the struct fields
// Returns:
//   - bool: true if reputation >= MinReputation, false otherwise
func (s *Student) CanDownload() bool {
	// Control flow: if/else to check reputation threshold
	if s.ReputationScore >= MinReputation {
		return true
	}
	return false
}

// UpdateReputation adjusts the student's reputation score
// Uses pointer receiver to modify the original struct
// Parameters:
//   - delta: Amount to add (positive) or subtract (negative)
func (s *Student) UpdateReputation(delta float64) {
	// Update the score
	s.ReputationScore += delta

	// Ensure score stays within valid range using control flow
	if s.ReputationScore > MaxReputation {
		s.ReputationScore = MaxReputation
	} else if s.ReputationScore < 0 {
		s.ReputationScore = 0
	}

	// Update leecher status based on upload/download ratio
	s.checkLeecherStatus()
}

// checkLeecherStatus determines if the student is a leecher
// A leecher downloads significantly more than they upload
func (s *Student) checkLeecherStatus() {
	// If downloads are more than 3x uploads, mark as leecher
	if s.TotalDownloads > 0 && s.TotalUploads > 0 {
		ratio := float64(s.TotalDownloads) / float64(s.TotalUploads)
		s.IsLeecher = ratio > 3.0
	} else if s.TotalDownloads > 5 && s.TotalUploads == 0 {
		// Downloaded files but never uploaded - definitely a leecher
		s.IsLeecher = true
	}
}

// RecordUpload increments upload count and rewards reputation
func (s *Student) RecordUpload() {
	s.TotalUploads++
	s.UpdateReputation(0.5) // Reward for contributing
	s.LastSeen = time.Now()
}

// RecordDownload increments download count
func (s *Student) RecordDownload() {
	s.TotalDownloads++
	s.LastSeen = time.Now()
	s.checkLeecherStatus()
}

// SetOnline updates the online status of the peer
func (s *Student) SetOnline(status bool) {
	s.IsOnline = status
	if status {
		s.LastSeen = time.Now()
	}
}

// GetAddress returns the full network address (IP:Port)
func (s *Student) GetAddress() string {
	return s.IPAddress + ":" + string(rune(s.Port))
}

// ToJSON converts the Student struct to JSON bytes
// Uses Go's encoding/json package for serialization
func (s *Student) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON populates the Student struct from JSON bytes
// Parameters:
//   - data: JSON byte slice containing student data
//
// Returns:
//   - error: nil if successful, error otherwise
func (s *Student) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}

// ============================================================================
// PEER REGISTRY - Map-based storage for all peers
// ============================================================================

// PeerRegistry stores all connected peers using a map
// Maps are key-value data structures in Go
// Key: Node ID (string), Value: Pointer to Student
type PeerRegistry struct {
	// peers is a map where string keys map to Student pointers
	peers map[string]*Student

	// mutex provides thread-safe access to the map
	mutex sync.RWMutex
}

// NewPeerRegistry creates a new empty peer registry
func NewPeerRegistry() *PeerRegistry {
	return &PeerRegistry{
		// Initialize an empty map using make()
		peers: make(map[string]*Student),
	}
}

// Register adds a new student to the registry
func (pr *PeerRegistry) Register(student *Student) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	// Add to map using student's ID as key
	pr.peers[student.ID] = student
}

// Unregister removes a student from the registry
func (pr *PeerRegistry) Unregister(id string) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	// Delete from map
	delete(pr.peers, id)
}

// Get retrieves a student by ID
// Returns:
//   - *Student: The student if found, nil otherwise
//   - bool: true if found, false otherwise
func (pr *PeerRegistry) Get(id string) (*Student, bool) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	// Map lookup with comma-ok idiom
	student, exists := pr.peers[id]
	return student, exists
}

// GetOnlinePeers returns a slice of all currently online peers
// Uses slice operations to build a dynamic list
func (pr *PeerRegistry) GetOnlinePeers() []*Student {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	// Create an empty slice to hold online peers
	var onlinePeers []*Student

	// Loop through all peers in the map
	for _, student := range pr.peers {
		if student.IsOnline {
			// Append to slice (dynamic array growth)
			onlinePeers = append(onlinePeers, student)
		}
	}

	return onlinePeers
}

// GetAllPeers returns all registered peers
func (pr *PeerRegistry) GetAllPeers() []*Student {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	peers := make([]*Student, 0, len(pr.peers))
	for _, student := range pr.peers {
		peers = append(peers, student)
	}
	return peers
}

// Count returns the total number of registered peers
func (pr *PeerRegistry) Count() int {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	return len(pr.peers)
}
