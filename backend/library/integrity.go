/*
================================================================================
INTEGRITY SERVICE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file handles file integrity verification using SHA-256 hashing.

Go Concepts Used:
- crypto/sha256: Cryptographic hashing
- io.Reader: Streaming interface for large files
- Error handling: Verification error management
================================================================================
*/

package library

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"knowledge-exchange/models"
)

// ============================================================================
// INTEGRITY SERVICE STRUCT
// ============================================================================

// IntegrityService provides file integrity verification capabilities
type IntegrityService struct {
	// Cache for previously verified files
	verifiedCache map[string]bool
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewIntegrityService creates a new IntegrityService
func NewIntegrityService() *IntegrityService {
	return &IntegrityService{
		verifiedCache: make(map[string]bool),
	}
}

// ============================================================================
// HASH GENERATION
// ============================================================================

// ComputeHash computes the SHA-256 hash of file content
// Parameters:
//   - content: File content as byte slice
//
// Returns:
//   - string: Hexadecimal hash string
func (is *IntegrityService) ComputeHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// ComputeFileHash computes the SHA-256 hash of a file
// Uses streaming to handle large files efficiently
// Parameters:
//   - filePath: Path to the file
//
// Returns:
//   - string: Hexadecimal hash string
//   - error: Error if file cannot be read
func (is *IntegrityService) ComputeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to compute hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ============================================================================
// VERIFICATION METHODS
// ============================================================================

// VerifyContent verifies that content matches an expected hash
// Parameters:
//   - content: The content to verify
//   - expectedHash: The expected hash value
//
// Returns:
//   - bool: true if content matches hash
func (is *IntegrityService) VerifyContent(content []byte, expectedHash string) bool {
	actualHash := is.ComputeHash(content)
	return actualHash == expectedHash
}

// VerifyFile verifies that a file matches an expected hash
// Parameters:
//   - filePath: Path to the file
//   - expectedHash: The expected hash value
//
// Returns:
//   - bool: true if file matches hash
//   - error: Error if file cannot be read
func (is *IntegrityService) VerifyFile(filePath, expectedHash string) (bool, error) {
	actualHash, err := is.ComputeFileHash(filePath)
	if err != nil {
		return false, err
	}
	return actualHash == expectedHash, nil
}

// VerifyAcademicFile verifies an AcademicFile against its stored checksum
// Parameters:
//   - file: The academic file record
//   - content: The file content
//
// Returns:
//   - bool: true if content matches stored checksum
func (is *IntegrityService) VerifyAcademicFile(file *models.AcademicFile, content []byte) bool {
	// Check cache first
	if verified, exists := is.verifiedCache[file.CID]; exists {
		return verified
	}

	// Verify against stored checksum
	result := is.VerifyContent(content, file.Checksum)

	// Cache the result
	is.verifiedCache[file.CID] = result

	return result
}

// ============================================================================
// CID VERIFICATION
// ============================================================================

// VerifyCID verifies that content matches a Content Identifier
// CID is derived from the content hash
func (is *IntegrityService) VerifyCID(content []byte, cid string) bool {
	// CID format: "kx-" + first 32 chars of hash
	actualHash := is.ComputeHash(content)
	expectedCID := "kx-" + actualHash[:32]
	return expectedCID == cid
}

// GenerateCID generates a Content Identifier from content
func (is *IntegrityService) GenerateCID(content []byte) string {
	hash := is.ComputeHash(content)
	return "kx-" + hash[:32]
}

// ============================================================================
// BATCH VERIFICATION
// ============================================================================

// VerifyResult holds the result of a verification operation
type VerifyResult struct {
	CID      string
	FilePath string
	Valid    bool
	Error    error
}

// VerifyMultipleFiles verifies multiple files
// Parameters:
//   - files: Map of CID to file path
//   - fileIndex: The file index for checksum lookup
//
// Returns:
//   - []VerifyResult: Results for each file
func (is *IntegrityService) VerifyMultipleFiles(files map[string]string, fileIndex *models.FileIndex) []VerifyResult {
	var results []VerifyResult

	for cid, filePath := range files {
		result := VerifyResult{
			CID:      cid,
			FilePath: filePath,
		}

		// Get the file record
		fileRecord, exists := fileIndex.Get(cid)
		if !exists {
			result.Valid = false
			result.Error = fmt.Errorf("file record not found for CID: %s", cid)
			results = append(results, result)
			continue
		}

		// Verify the file
		valid, err := is.VerifyFile(filePath, fileRecord.Checksum)
		result.Valid = valid
		result.Error = err

		results = append(results, result)
	}

	return results
}

// ============================================================================
// CACHE MANAGEMENT
// ============================================================================

// ClearCache clears the verification cache
func (is *IntegrityService) ClearCache() {
	is.verifiedCache = make(map[string]bool)
}

// InvalidateCache removes a specific entry from the cache
func (is *IntegrityService) InvalidateCache(cid string) {
	delete(is.verifiedCache, cid)
}

// GetCacheSize returns the number of cached verifications
func (is *IntegrityService) GetCacheSize() int {
	return len(is.verifiedCache)
}
