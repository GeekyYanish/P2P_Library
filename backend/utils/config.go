/*
================================================================================
CONFIGURATION - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file contains system configuration constants and settings.

Go Concepts Used:
- Constants: Fixed system values
- Custom types: Type aliases for clarity
- Structs: Configuration grouping
================================================================================
*/

package utils

import (
	"encoding/json"
	"os"
	"time"
)

// ============================================================================
// SYSTEM CONSTANTS
// ============================================================================

const (
	// Application Info
	AppName    = "Knowledge Exchange"
	AppVersion = "1.0.0"

	// Network Defaults
	DefaultServerPort = 8080
	DefaultAPIPort    = 3000

	// Reputation System
	MinReputationToDownload = 3.0
	MaxReputation           = 10.0
	DefaultReputation       = 5.0
	ReputationGainPerUpload = 0.5
	ReputationLossPerLeech  = 0.2

	// Throttling
	LeecherBandwidthLimit = 100 * 1024       // 100 KB/s for leechers
	NormalBandwidthLimit  = 1024 * 1024      // 1 MB/s for normal users
	PremiumBandwidthLimit = 10 * 1024 * 1024 // 10 MB/s for high reputation

	// File Limits
	MaxFileSizeBytes       = 100 * 1024 * 1024 // 100 MB
	MaxConcurrentUploads   = 5
	MaxConcurrentDownloads = 3

	// Timeouts
	PeerTimeoutSeconds       = 30
	TransferTimeoutSeconds   = 300
	HeartbeatIntervalSeconds = 10

	// Storage
	DefaultDataDir = "./data"
	SharedFilesDir = "./data/sharedFiles"
	TempDir        = "./data/temp"
)

// ============================================================================
// CONFIGURATION STRUCT
// ============================================================================

// Config holds all application configuration
type Config struct {
	// Server Settings
	ServerPort int    `json:"server_port"`
	APIPort    int    `json:"api_port"`
	HostIP     string `json:"host_ip"`

	// Peer Identity
	PeerID   string `json:"peer_id"`
	PeerName string `json:"peer_name"`

	// Storage Paths
	DataDir        string `json:"data_dir"`
	SharedFilesDir string `json:"shared_files_dir"`
	TempDir        string `json:"temp_dir"`

	// Reputation Settings
	MinReputation float64 `json:"min_reputation"`
	MaxReputation float64 `json:"max_reputation"`

	// Throttling Settings
	LeecherBandwidth int64 `json:"leecher_bandwidth"`
	NormalBandwidth  int64 `json:"normal_bandwidth"`

	// Timeouts
	PeerTimeout     time.Duration `json:"peer_timeout"`
	TransferTimeout time.Duration `json:"transfer_timeout"`

	// Limits
	MaxFileSize     int64 `json:"max_file_size"`
	MaxConcurrentTx int   `json:"max_concurrent_tx"`

	// Feature Flags
	EnableThrottling bool `json:"enable_throttling"`
	EnableRatings    bool `json:"enable_ratings"`
	EnableEncryption bool `json:"enable_encryption"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		ServerPort:       DefaultServerPort,
		APIPort:          DefaultAPIPort,
		HostIP:           "127.0.0.1",
		PeerID:           "",
		PeerName:         "Anonymous Peer",
		DataDir:          DefaultDataDir,
		SharedFilesDir:   SharedFilesDir,
		TempDir:          TempDir,
		MinReputation:    MinReputationToDownload,
		MaxReputation:    MaxReputation,
		LeecherBandwidth: LeecherBandwidthLimit,
		NormalBandwidth:  NormalBandwidthLimit,
		PeerTimeout:      time.Duration(PeerTimeoutSeconds) * time.Second,
		TransferTimeout:  time.Duration(TransferTimeoutSeconds) * time.Second,
		MaxFileSize:      MaxFileSizeBytes,
		MaxConcurrentTx:  MaxConcurrentDownloads,
		EnableThrottling: true,
		EnableRatings:    true,
		EnableEncryption: false,
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves configuration to a JSON file
func (c *Config) SaveConfig(filePath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// ============================================================================
// ALLOWED FILE TYPES
// ============================================================================

// AllowedFileTypes returns the list of permitted file extensions
func AllowedFileTypes() []string {
	return []string{
		".pdf",
		".doc",
		".docx",
		".ppt",
		".pptx",
		".txt",
		".md",
		".epub",
		".odt",
		".xls",
		".xlsx",
	}
}

// IsAllowedFileType checks if a file extension is allowed
func IsAllowedFileType(ext string) bool {
	for _, allowed := range AllowedFileTypes() {
		if ext == allowed {
			return true
		}
	}
	return false
}

// ============================================================================
// DIRECTORY SETUP
// ============================================================================

// EnsureDirectories creates necessary directories
func EnsureDirectories() error {
	dirs := []string{
		DefaultDataDir,
		SharedFilesDir,
		TempDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
