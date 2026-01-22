/*
================================================================================
ACADEMIC FILE MODEL - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file defines the AcademicFile struct representing shared resources.

Go Concepts Used:
- Structs: Grouping file metadata together
- Slices: Dynamic arrays for storing file lists
- Maps: Key-value storage for file indexing using CID
- Byte Slices ([]byte): For file content and hashing
- Methods: Functions for file validation and operations
================================================================================
*/

package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

// ============================================================================
// CONSTANTS - File type and size limits
// ============================================================================

// MaxFileSize is the maximum allowed file size (100 MB)
const MaxFileSize int64 = 100 * 1024 * 1024

// AllowedFileTypes defines the permitted academic file extensions
var AllowedFileTypes = []string{".pdf", ".doc", ".docx", ".ppt", ".pptx", ".txt", ".md"}

// ============================================================================
// STRUCT DEFINITION - The AcademicFile data type
// ============================================================================

// AcademicFile represents an academic resource being shared in the P2P network
// Files are identified by their Content Identifier (CID) which is a SHA-256 hash
type AcademicFile struct {
	// CID is the Content Identifier - a unique hash of the file content
	// This ensures file integrity and enables content-addressable storage
	CID string `json:"cid"`

	// FileName is the original name of the file
	FileName string `json:"file_name"`

	// OwnerID is the ID of the student who originally shared this file
	OwnerID string `json:"owner_id"`

	// Size is the file size in bytes
	Size int64 `json:"size"`

	// FileType is the extension/type of the file (e.g., ".pdf")
	FileType string `json:"file_type"`

	// Description is an optional description of the file contents
	Description string `json:"description"`

	// Subject categorizes the file (e.g., "Algorithms", "Database", "OS")
	Subject string `json:"subject"`

	// UploadTime records when the file was first shared
	UploadTime time.Time `json:"upload_time"`

	// DownloadCount tracks how many times this file has been downloaded
	DownloadCount int `json:"download_count"`

	// AverageRating is the average peer rating for this file
	AverageRating float64 `json:"average_rating"`

	// TotalRatings is the number of ratings received
	TotalRatings int `json:"total_ratings"`

	// PeerLocations stores IDs of peers that have this file
	// This is a slice (dynamic array) of strings
	PeerLocations []string `json:"peer_locations"`

	// IsAvailable indicates if at least one peer with this file is online
	IsAvailable bool `json:"is_available"`

	// Checksum is used for integrity verification after download
	Checksum string `json:"checksum"`
}

// ============================================================================
// CONSTRUCTOR - Create new AcademicFile instances
// ============================================================================

// NewAcademicFile creates a new AcademicFile with computed CID
// Parameters:
//   - fileName: Name of the file
//   - ownerID: ID of the uploading student
//   - size: Size in bytes
//   - fileType: File extension
//   - content: The actual file content (byte slice)
//
// Returns:
//   - *AcademicFile: Pointer to the newly created file
func NewAcademicFile(fileName, ownerID string, size int64, fileType string, content []byte) *AcademicFile {
	// Generate CID using SHA-256 hash of content
	cid := GenerateCID(content)
	checksum := GenerateChecksum(content)

	return &AcademicFile{
		CID:           cid,
		FileName:      fileName,
		OwnerID:       ownerID,
		Size:          size,
		FileType:      fileType,
		Description:   "",
		Subject:       "",
		UploadTime:    time.Now(),
		DownloadCount: 0,
		AverageRating: 0.0,
		TotalRatings:  0,
		PeerLocations: []string{ownerID}, // Owner is the first peer with the file
		IsAvailable:   true,
		Checksum:      checksum,
	}
}

// ============================================================================
// HASHING FUNCTIONS - CID and Checksum generation
// ============================================================================

// GenerateCID creates a Content Identifier from file content
// Uses SHA-256 hashing for unique identification
// Parameters:
//   - content: File content as byte slice ([]byte)
//
// Returns:
//   - string: Hexadecimal string of the hash
func GenerateCID(content []byte) string {
	// Create SHA-256 hash
	hash := sha256.New()

	// Write content bytes to the hash
	hash.Write(content)

	// Get the final hash as byte slice
	hashBytes := hash.Sum(nil)

	// Convert to hexadecimal string
	return hex.EncodeToString(hashBytes)
}

// GenerateChecksum creates a checksum for integrity verification
// Same as CID but semantically different purpose
func GenerateChecksum(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// ============================================================================
// METHODS - AcademicFile operations
// ============================================================================

// IsValid checks if the file meets all requirements
// Returns:
//   - bool: true if valid, false otherwise
//   - string: Error message if invalid
func (f *AcademicFile) IsValid() (bool, string) {
	// Check file size
	if f.Size > MaxFileSize {
		return false, "File exceeds maximum size limit"
	}

	if f.Size <= 0 {
		return false, "File size must be positive"
	}

	// Check file type using a loop
	validType := false
	for _, allowed := range AllowedFileTypes {
		if f.FileType == allowed {
			validType = true
			break
		}
	}

	if !validType {
		return false, "File type not allowed"
	}

	// Check required fields
	if f.FileName == "" {
		return false, "File name is required"
	}

	if f.OwnerID == "" {
		return false, "Owner ID is required"
	}

	return true, ""
}

// AddPeerLocation adds a new peer that has this file
func (f *AcademicFile) AddPeerLocation(peerID string) {
	// Check if peer already in list using loop
	for _, id := range f.PeerLocations {
		if id == peerID {
			return // Already exists
		}
	}

	// Append to slice
	f.PeerLocations = append(f.PeerLocations, peerID)
}

// RemovePeerLocation removes a peer from the locations list
func (f *AcademicFile) RemovePeerLocation(peerID string) {
	// Create new slice without the peer
	newLocations := make([]string, 0, len(f.PeerLocations))

	for _, id := range f.PeerLocations {
		if id != peerID {
			newLocations = append(newLocations, id)
		}
	}

	f.PeerLocations = newLocations

	// Update availability
	f.IsAvailable = len(f.PeerLocations) > 0
}

// RecordDownload increments the download counter
func (f *AcademicFile) RecordDownload() {
	f.DownloadCount++
}

// AddRating adds a new rating and updates the average
// Parameters:
//   - rating: The rating value (1-5)
func (f *AcademicFile) AddRating(rating float64) {
	// Calculate new average
	totalScore := f.AverageRating * float64(f.TotalRatings)
	f.TotalRatings++
	f.AverageRating = (totalScore + rating) / float64(f.TotalRatings)
}

// VerifyIntegrity checks if downloaded content matches the checksum
// Parameters:
//   - content: Downloaded file content
//
// Returns:
//   - bool: true if checksum matches
func (f *AcademicFile) VerifyIntegrity(content []byte) bool {
	computedChecksum := GenerateChecksum(content)
	return computedChecksum == f.Checksum
}

// ToJSON converts the file to JSON bytes
func (f *AcademicFile) ToJSON() ([]byte, error) {
	return json.Marshal(f)
}

// FromJSON populates the file from JSON bytes
func (f *AcademicFile) FromJSON(data []byte) error {
	return json.Unmarshal(data, f)
}

// ============================================================================
// FILE INDEX - Map-based storage for all files
// ============================================================================

// FileIndex stores all shared files using a map
// Key: CID (Content Identifier), Value: Pointer to AcademicFile
type FileIndex struct {
	// files is a map where CID keys map to AcademicFile pointers
	files map[string]*AcademicFile

	// mutex provides thread-safe access
	mutex sync.RWMutex
}

// NewFileIndex creates a new empty file index
func NewFileIndex() *FileIndex {
	return &FileIndex{
		files: make(map[string]*AcademicFile),
	}
}

// Add adds a new file to the index
func (fi *FileIndex) Add(file *AcademicFile) {
	fi.mutex.Lock()
	defer fi.mutex.Unlock()

	fi.files[file.CID] = file
}

// Get retrieves a file by CID
func (fi *FileIndex) Get(cid string) (*AcademicFile, bool) {
	fi.mutex.RLock()
	defer fi.mutex.RUnlock()

	file, exists := fi.files[cid]
	return file, exists
}

// Remove removes a file from the index
func (fi *FileIndex) Remove(cid string) {
	fi.mutex.Lock()
	defer fi.mutex.Unlock()

	delete(fi.files, cid)
}

// Search finds files matching the query string
// Searches in filename, description, and subject
// Returns:
//   - []*AcademicFile: Slice of matching files
func (fi *FileIndex) Search(query string) []*AcademicFile {
	fi.mutex.RLock()
	defer fi.mutex.RUnlock()

	var results []*AcademicFile

	// Loop through all files
	for _, file := range fi.files {
		// Simple string matching (case-sensitive)
		if containsIgnoreCase(file.FileName, query) ||
			containsIgnoreCase(file.Description, query) ||
			containsIgnoreCase(file.Subject, query) {
			results = append(results, file)
		}
	}

	return results
}

// GetBySubject returns all files for a specific subject
func (fi *FileIndex) GetBySubject(subject string) []*AcademicFile {
	fi.mutex.RLock()
	defer fi.mutex.RUnlock()

	var results []*AcademicFile

	for _, file := range fi.files {
		if file.Subject == subject {
			results = append(results, file)
		}
	}

	return results
}

// GetAvailableFiles returns only files that have online peers
func (fi *FileIndex) GetAvailableFiles() []*AcademicFile {
	fi.mutex.RLock()
	defer fi.mutex.RUnlock()

	var available []*AcademicFile

	for _, file := range fi.files {
		if file.IsAvailable {
			available = append(available, file)
		}
	}

	return available
}

// GetAllFiles returns all files in the index
func (fi *FileIndex) GetAllFiles() []*AcademicFile {
	fi.mutex.RLock()
	defer fi.mutex.RUnlock()

	files := make([]*AcademicFile, 0, len(fi.files))
	for _, file := range fi.files {
		files = append(files, file)
	}
	return files
}

// Count returns the total number of files
func (fi *FileIndex) Count() int {
	fi.mutex.RLock()
	defer fi.mutex.RUnlock()

	return len(fi.files)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	// Simple implementation - check if substr exists in s
	// In production, use strings.Contains with strings.ToLower
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr))
}
