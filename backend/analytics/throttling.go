/*
================================================================================
THROTTLING SERVICE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements bandwidth throttling for fair resource usage.

Go Concepts Used:
- Goroutines: Managing throttled transfers
- Channels: Rate limiting token buckets
- time.Ticker: Periodic token replenishment
- Interfaces: ThrottledConnection abstraction
================================================================================
*/

package analytics

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// Bandwidth limits in bytes per second
	LeecherBandwidth = 50 * 1024       // 50 KB/s for leechers
	NormalBandwidth  = 500 * 1024      // 500 KB/s for normal users
	PremiumBandwidth = 5 * 1024 * 1024 // 5 MB/s for high-reputation users

	// Reputation thresholds for bandwidth tiers
	LeecherThreshold = 3.0
	PremiumThreshold = 8.0

	// Token bucket parameters
	TokenBucketSize = 10
	RefillInterval  = 100 * time.Millisecond
)

// ============================================================================
// BANDWIDTH TIER
// ============================================================================

// BandwidthTier represents a bandwidth allocation tier
type BandwidthTier int

const (
	TierLeecher BandwidthTier = iota
	TierNormal
	TierPremium
)

// String returns the tier name
func (t BandwidthTier) String() string {
	switch t {
	case TierLeecher:
		return "Leecher"
	case TierNormal:
		return "Normal"
	case TierPremium:
		return "Premium"
	default:
		return "Unknown"
	}
}

// GetBandwidth returns the bandwidth limit for a tier
func (t BandwidthTier) GetBandwidth() int64 {
	switch t {
	case TierLeecher:
		return LeecherBandwidth
	case TierNormal:
		return NormalBandwidth
	case TierPremium:
		return PremiumBandwidth
	default:
		return NormalBandwidth
	}
}

// ============================================================================
// THROTTLER STRUCT
// ============================================================================

// Throttler manages bandwidth allocation for a peer
type Throttler struct {
	peerID    string
	tier      BandwidthTier
	bandwidth int64 // bytes per second
	tokens    int64 // current available tokens
	maxTokens int64 // maximum tokens (bucket size)
	tokenSize int64 // bytes per token
	mutex     sync.Mutex
	ticker    *time.Ticker
	stopChan  chan struct{}
	isRunning bool
}

// NewThrottler creates a new throttler for a peer
func NewThrottler(peerID string, reputation float64) *Throttler {
	tier := determineTier(reputation)
	bandwidth := tier.GetBandwidth()

	// Calculate token size (how many bytes per token)
	tokenSize := bandwidth / TokenBucketSize

	t := &Throttler{
		peerID:    peerID,
		tier:      tier,
		bandwidth: bandwidth,
		tokens:    TokenBucketSize, // Start with full bucket
		maxTokens: TokenBucketSize,
		tokenSize: tokenSize,
		stopChan:  make(chan struct{}),
		isRunning: false,
	}

	return t
}

// determineTier determines bandwidth tier based on reputation
func determineTier(reputation float64) BandwidthTier {
	if reputation < LeecherThreshold {
		return TierLeecher
	} else if reputation >= PremiumThreshold {
		return TierPremium
	}
	return TierNormal
}

// ============================================================================
// THROTTLER LIFECYCLE
// ============================================================================

// Start begins the token refill goroutine
func (t *Throttler) Start() {
	if t.isRunning {
		return
	}

	t.isRunning = true
	t.ticker = time.NewTicker(RefillInterval)

	// Token refill goroutine
	go func() {
		tokensPerRefill := int64(1) // Add 1 token per interval

		for {
			select {
			case <-t.ticker.C:
				t.mutex.Lock()
				if t.tokens < t.maxTokens {
					t.tokens += tokensPerRefill
					if t.tokens > t.maxTokens {
						t.tokens = t.maxTokens
					}
				}
				t.mutex.Unlock()

			case <-t.stopChan:
				t.ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the throttler
func (t *Throttler) Stop() {
	if t.isRunning {
		t.isRunning = false
		close(t.stopChan)
	}
}

// ============================================================================
// THROTTLING METHODS
// ============================================================================

// Acquire acquires tokens for a given number of bytes
// Blocks until tokens are available
// Parameters:
//   - bytes: Number of bytes to acquire tokens for
//
// Returns:
//   - int64: Number of bytes actually allowed
func (t *Throttler) Acquire(bytes int64) int64 {
	// Calculate tokens needed
	tokensNeeded := (bytes + t.tokenSize - 1) / t.tokenSize // Round up

	// Wait for tokens to be available
	for {
		t.mutex.Lock()
		if t.tokens >= tokensNeeded {
			t.tokens -= tokensNeeded
			t.mutex.Unlock()
			return bytes
		} else if t.tokens > 0 {
			// Use available tokens for partial transfer
			allowedBytes := t.tokens * t.tokenSize
			t.tokens = 0
			t.mutex.Unlock()
			return allowedBytes
		}
		t.mutex.Unlock()

		// Wait for token refill
		time.Sleep(RefillInterval)
	}
}

// GetAvailableBytes returns currently available bandwidth
func (t *Throttler) GetAvailableBytes() int64 {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.tokens * t.tokenSize
}

// GetTier returns the current bandwidth tier
func (t *Throttler) GetTier() BandwidthTier {
	return t.tier
}

// GetBandwidth returns the bandwidth limit
func (t *Throttler) GetBandwidth() int64 {
	return t.bandwidth
}

// UpdateReputation updates the throttler based on new reputation
func (t *Throttler) UpdateReputation(newReputation float64) {
	newTier := determineTier(newReputation)

	if newTier != t.tier {
		t.mutex.Lock()
		t.tier = newTier
		t.bandwidth = newTier.GetBandwidth()
		t.tokenSize = t.bandwidth / TokenBucketSize
		t.mutex.Unlock()
	}
}

// ============================================================================
// THROTTLED READER/WRITER
// ============================================================================

// ThrottledReader wraps an io.Reader with throttling
type ThrottledReader struct {
	reader    io.Reader
	throttler *Throttler
}

// NewThrottledReader creates a new throttled reader
func NewThrottledReader(reader io.Reader, throttler *Throttler) *ThrottledReader {
	return &ThrottledReader{
		reader:    reader,
		throttler: throttler,
	}
}

// Read implements io.Reader with throttling
func (tr *ThrottledReader) Read(p []byte) (n int, err error) {
	// Calculate allowed bytes
	allowed := tr.throttler.Acquire(int64(len(p)))

	// Read up to allowed bytes
	return tr.reader.Read(p[:allowed])
}

// ThrottledWriter wraps an io.Writer with throttling
type ThrottledWriter struct {
	writer    io.Writer
	throttler *Throttler
}

// NewThrottledWriter creates a new throttled writer
func NewThrottledWriter(writer io.Writer, throttler *Throttler) *ThrottledWriter {
	return &ThrottledWriter{
		writer:    writer,
		throttler: throttler,
	}
}

// Write implements io.Writer with throttling
func (tw *ThrottledWriter) Write(p []byte) (n int, err error) {
	written := 0
	for written < len(p) {
		// Acquire tokens for remaining bytes
		remaining := int64(len(p) - written)
		allowed := tw.throttler.Acquire(remaining)

		// Write allowed bytes
		n, err := tw.writer.Write(p[written : written+int(allowed)])
		if err != nil {
			return written + n, err
		}
		written += n
	}
	return written, nil
}

// ============================================================================
// THROTTLING MANAGER
// ============================================================================

// ThrottlingManager manages throttlers for all peers
type ThrottlingManager struct {
	throttlers map[string]*Throttler
	mutex      sync.RWMutex
	enabled    bool
}

// NewThrottlingManager creates a new throttling manager
func NewThrottlingManager() *ThrottlingManager {
	return &ThrottlingManager{
		throttlers: make(map[string]*Throttler),
		enabled:    true,
	}
}

// GetThrottler gets or creates a throttler for a peer
func (tm *ThrottlingManager) GetThrottler(peerID string, reputation float64) *Throttler {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if throttler, exists := tm.throttlers[peerID]; exists {
		throttler.UpdateReputation(reputation)
		return throttler
	}

	throttler := NewThrottler(peerID, reputation)
	throttler.Start()
	tm.throttlers[peerID] = throttler
	return throttler
}

// RemoveThrottler removes and stops a throttler
func (tm *ThrottlingManager) RemoveThrottler(peerID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if throttler, exists := tm.throttlers[peerID]; exists {
		throttler.Stop()
		delete(tm.throttlers, peerID)
	}
}

// SetEnabled enables or disables throttling globally
func (tm *ThrottlingManager) SetEnabled(enabled bool) {
	tm.enabled = enabled
}

// IsEnabled returns whether throttling is enabled
func (tm *ThrottlingManager) IsEnabled() bool {
	return tm.enabled
}

// GetStats returns throttling statistics
func (tm *ThrottlingManager) GetStats() map[string]interface{} {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	leecherCount := 0
	normalCount := 0
	premiumCount := 0

	for _, t := range tm.throttlers {
		switch t.tier {
		case TierLeecher:
			leecherCount++
		case TierNormal:
			normalCount++
		case TierPremium:
			premiumCount++
		}
	}

	return map[string]interface{}{
		"enabled":       tm.enabled,
		"total_peers":   len(tm.throttlers),
		"leecher_count": leecherCount,
		"normal_count":  normalCount,
		"premium_count": premiumCount,
	}
}

// GetPeerInfo returns throttling info for a specific peer
func (tm *ThrottlingManager) GetPeerInfo(peerID string) (map[string]interface{}, error) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	throttler, exists := tm.throttlers[peerID]
	if !exists {
		return nil, fmt.Errorf("no throttler for peer: %s", peerID)
	}

	return map[string]interface{}{
		"peer_id":         peerID,
		"tier":            throttler.tier.String(),
		"bandwidth_limit": throttler.bandwidth,
		"available_bytes": throttler.GetAvailableBytes(),
	}, nil
}

// StopAll stops all throttlers
func (tm *ThrottlingManager) StopAll() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	for _, t := range tm.throttlers {
		t.Stop()
	}
	tm.throttlers = make(map[string]*Throttler)
}
