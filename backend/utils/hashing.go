/*
================================================================================
HASHING UTILITIES - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file provides cryptographic hashing functions for file integrity.

Go Concepts Used:
- crypto/sha256: Cryptographic hashing
- Byte Slices: For binary data handling
- Functions: Modular utility functions
================================================================================
*/

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// ============================================================================
// HASHING FUNCTIONS
// ============================================================================

// HashBytes computes the SHA-256 hash of a byte slice
// Parameters:
//   - data: The byte slice to hash
//
// Returns:
//   - string: Hexadecimal string of the hash
func HashBytes(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// HashString computes the SHA-256 hash of a string
func HashString(s string) string {
	return HashBytes([]byte(s))
}

// HashFile computes the SHA-256 hash of a file
// Parameters:
//   - filePath: Path to the file to hash
//
// Returns:
//   - string: Hexadecimal string of the hash
//   - error: Error if file cannot be read
func HashFile(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create hash instance
	hash := sha256.New()

	// Copy file contents to hash (efficient for large files)
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	// Return hex-encoded hash
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GenerateCID generates a Content Identifier for academic content
// CIDs are unique identifiers based on content hash
func GenerateCID(content []byte) string {
	// Prefix with "kx-" for Knowledge Exchange
	hash := HashBytes(content)
	return "kx-" + hash[:32] // Use first 32 chars for shorter CID
}

// VerifyHash checks if content matches an expected hash
// Parameters:
//   - content: The content to verify
//   - expectedHash: The hash to compare against
//
// Returns:
//   - bool: true if hashes match
func VerifyHash(content []byte, expectedHash string) bool {
	computedHash := HashBytes(content)
	return computedHash == expectedHash
}

// VerifyFileHash checks if a file matches an expected hash
func VerifyFileHash(filePath, expectedHash string) (bool, error) {
	actualHash, err := HashFile(filePath)
	if err != nil {
		return false, err
	}
	return actualHash == expectedHash, nil
}

// GeneratePeerID generates a unique peer ID from identifying information
func GeneratePeerID(name, ip string, port int) string {
	data := fmt.Sprintf("%s:%s:%d:%d", name, ip, port, getCurrentTimestamp())
	hash := HashString(data)
	return "peer-" + hash[:16]
}

// getCurrentTimestamp returns the current Unix timestamp
func getCurrentTimestamp() int64 {
	return 0 // Simplified - would use time.Now().Unix() in production
}

// ============================================================================
// CHECKSUM FUNCTIONS
// ============================================================================

// ComputeChecksum generates a checksum for data integrity
func ComputeChecksum(data []byte) string {
	return HashBytes(data)
}

// VerifyChecksum verifies data against a checksum
func VerifyChecksum(data []byte, checksum string) bool {
	return ComputeChecksum(data) == checksum
}

// ChunkHash splits data into chunks and hashes each
// Useful for parallel hashing of large files
func ChunkHash(data []byte, chunkSize int) []string {
	if chunkSize <= 0 {
		chunkSize = 1024 * 1024 // Default 1MB chunks
	}

	var hashes []string
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		hashes = append(hashes, HashBytes(data[i:end]))
	}
	return hashes
}
