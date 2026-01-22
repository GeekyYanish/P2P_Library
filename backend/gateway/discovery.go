/*
================================================================================
DISCOVERY SERVICE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file implements peer discovery for the P2P network.

Go Concepts Used:
- Goroutines: Background peer checking
- Channels: Discovery event broadcasting
- Maps: Peer tracking
- Time: Heartbeat and timeout handling
================================================================================
*/

package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"knowledge-exchange/models"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// Discovery intervals
	HeartbeatInterval = 10 * time.Second
	PeerTimeout       = 30 * time.Second
	CleanupInterval   = 1 * time.Minute

	// Discovery message types
	DiscoveryAnnounce = "ANNOUNCE"
	DiscoveryPing     = "PING"
	DiscoveryPong     = "PONG"
	DiscoveryLeave    = "LEAVE"
)

// ============================================================================
// DISCOVERY MESSAGE
// ============================================================================

// DiscoveryMessage represents a peer discovery message
type DiscoveryMessage struct {
	Type      string    `json:"type"`
	PeerID    string    `json:"peer_id"`
	PeerName  string    `json:"peer_name"`
	Address   string    `json:"address"`
	Port      int       `json:"port"`
	Timestamp time.Time `json:"timestamp"`
}

// ============================================================================
// DISCOVERY SERVICE STRUCT
// ============================================================================

// Discovery handles peer discovery and health monitoring
type Discovery struct {
	// Peer registry
	peerRegistry *models.PeerRegistry

	// Known peer addresses
	knownPeers map[string]time.Time // peerID -> last seen

	// Event subscribers
	subscribers []chan DiscoveryEvent

	// Mutex for thread-safe operations
	mutex sync.RWMutex

	// Control channels
	stopChan  chan struct{}
	eventChan chan DiscoveryEvent

	// State
	isRunning bool
	localPeer *models.Student
}

// DiscoveryEvent represents a discovery event
type DiscoveryEvent struct {
	Type    string          `json:"type"`
	PeerID  string          `json:"peer_id"`
	Peer    *models.Student `json:"peer,omitempty"`
	Message string          `json:"message,omitempty"`
}

// Event types
const (
	EventPeerJoined  = "PEER_JOINED"
	EventPeerLeft    = "PEER_LEFT"
	EventPeerTimeout = "PEER_TIMEOUT"
	EventPeerUpdated = "PEER_UPDATED"
)

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewDiscovery creates a new Discovery service
func NewDiscovery(peerRegistry *models.PeerRegistry) *Discovery {
	return &Discovery{
		peerRegistry: peerRegistry,
		knownPeers:   make(map[string]time.Time),
		subscribers:  make([]chan DiscoveryEvent, 0),
		stopChan:     make(chan struct{}),
		eventChan:    make(chan DiscoveryEvent, 100),
		isRunning:    false,
	}
}

// ============================================================================
// SERVICE LIFECYCLE
// ============================================================================

// Start begins the discovery service
func (d *Discovery) Start() {
	if d.isRunning {
		return
	}

	d.isRunning = true

	// Start event broadcaster
	go d.broadcastEvents()

	// Start heartbeat sender
	go d.sendHeartbeats()

	// Start cleanup routine
	go d.cleanupStale()

	log.Println("Discovery service started")
}

// Stop stops the discovery service
func (d *Discovery) Stop() {
	if !d.isRunning {
		return
	}

	d.isRunning = false
	close(d.stopChan)
	close(d.eventChan)

	// Close all subscriber channels
	for _, ch := range d.subscribers {
		close(ch)
	}

	log.Println("Discovery service stopped")
}

// SetLocalPeer sets the local peer identity
func (d *Discovery) SetLocalPeer(peer *models.Student) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.localPeer = peer
}

// ============================================================================
// PEER DISCOVERY
// ============================================================================

// AnnouncePeer announces local peer to the network
func (d *Discovery) AnnouncePeer(peer *models.Student) error {
	msg := &DiscoveryMessage{
		Type:      DiscoveryAnnounce,
		PeerID:    peer.ID,
		PeerName:  peer.Name,
		Address:   peer.IPAddress,
		Port:      peer.Port,
		Timestamp: time.Now(),
	}

	return d.broadcastMessage(msg)
}

// RegisterPeer registers a discovered peer
func (d *Discovery) RegisterPeer(msg *DiscoveryMessage) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if already known
	_, exists := d.knownPeers[msg.PeerID]

	// Update last seen
	d.knownPeers[msg.PeerID] = time.Now()

	// Get or create peer in registry
	peer, peerExists := d.peerRegistry.Get(msg.PeerID)

	if !peerExists {
		// Create new peer
		peer = models.NewStudent(msg.PeerID, msg.PeerName, msg.Address, msg.Port)
		d.peerRegistry.Register(peer)

		// Emit join event
		d.eventChan <- DiscoveryEvent{
			Type:   EventPeerJoined,
			PeerID: msg.PeerID,
			Peer:   peer,
		}
	} else if !exists {
		// Peer returned online
		peer.SetOnline(true)
		d.eventChan <- DiscoveryEvent{
			Type:   EventPeerUpdated,
			PeerID: msg.PeerID,
			Peer:   peer,
		}
	}
}

// HandleLeave handles a peer leaving the network
func (d *Discovery) HandleLeave(peerID string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.knownPeers, peerID)

	if peer, exists := d.peerRegistry.Get(peerID); exists {
		peer.SetOnline(false)
		d.eventChan <- DiscoveryEvent{
			Type:   EventPeerLeft,
			PeerID: peerID,
			Peer:   peer,
		}
	}
}

// ============================================================================
// HEARTBEAT
// ============================================================================

// sendHeartbeats periodically sends heartbeats to known peers
func (d *Discovery) sendHeartbeats() {
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.pingAllPeers()
		case <-d.stopChan:
			return
		}
	}
}

// pingAllPeers sends ping to all known peers
func (d *Discovery) pingAllPeers() {
	d.mutex.RLock()
	peerIDs := make([]string, 0, len(d.knownPeers))
	for id := range d.knownPeers {
		peerIDs = append(peerIDs, id)
	}
	d.mutex.RUnlock()

	for _, peerID := range peerIDs {
		go d.pingPeer(peerID)
	}
}

// pingPeer pings a specific peer
func (d *Discovery) pingPeer(peerID string) {
	peer, exists := d.peerRegistry.Get(peerID)
	if !exists {
		return
	}

	address := fmt.Sprintf("%s:%d", peer.IPAddress, peer.Port)

	// Try to connect
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		// Peer not responding
		return
	}
	defer conn.Close()

	// Send ping
	msg := &DiscoveryMessage{
		Type:      DiscoveryPing,
		PeerID:    d.getLocalPeerID(),
		Timestamp: time.Now(),
	}

	data, _ := json.Marshal(msg)
	conn.Write(data)

	// Update last seen on successful ping
	d.mutex.Lock()
	d.knownPeers[peerID] = time.Now()
	d.mutex.Unlock()
}

// getLocalPeerID returns the local peer ID
func (d *Discovery) getLocalPeerID() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.localPeer != nil {
		return d.localPeer.ID
	}
	return "unknown"
}

// ============================================================================
// CLEANUP
// ============================================================================

// cleanupStale removes stale peers
func (d *Discovery) cleanupStale() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.removeStale()
		case <-d.stopChan:
			return
		}
	}
}

// removeStale removes peers that haven't been seen recently
func (d *Discovery) removeStale() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	now := time.Now()
	stale := make([]string, 0)

	for peerID, lastSeen := range d.knownPeers {
		if now.Sub(lastSeen) > PeerTimeout {
			stale = append(stale, peerID)
		}
	}

	for _, peerID := range stale {
		delete(d.knownPeers, peerID)

		if peer, exists := d.peerRegistry.Get(peerID); exists {
			peer.SetOnline(false)
			d.eventChan <- DiscoveryEvent{
				Type:    EventPeerTimeout,
				PeerID:  peerID,
				Peer:    peer,
				Message: "Peer timed out",
			}
		}
	}
}

// ============================================================================
// EVENT BROADCASTING
// ============================================================================

// broadcastEvents broadcasts events to all subscribers
func (d *Discovery) broadcastEvents() {
	for {
		select {
		case event, ok := <-d.eventChan:
			if !ok {
				return
			}
			d.mutex.RLock()
			for _, ch := range d.subscribers {
				select {
				case ch <- event:
				default:
					// Skip if subscriber is not ready
				}
			}
			d.mutex.RUnlock()
		case <-d.stopChan:
			return
		}
	}
}

// Subscribe returns a channel for discovery events
func (d *Discovery) Subscribe() <-chan DiscoveryEvent {
	ch := make(chan DiscoveryEvent, 10)
	d.mutex.Lock()
	d.subscribers = append(d.subscribers, ch)
	d.mutex.Unlock()
	return ch
}

// ============================================================================
// MESSAGE BROADCASTING
// ============================================================================

// broadcastMessage broadcasts a message to all known peers
func (d *Discovery) broadcastMessage(msg *DiscoveryMessage) error {
	d.mutex.RLock()
	peers := d.peerRegistry.GetOnlinePeers()
	d.mutex.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, peer := range peers {
		go func(p *models.Student) {
			address := fmt.Sprintf("%s:%d", p.IPAddress, p.Port)
			conn, err := net.DialTimeout("tcp", address, 5*time.Second)
			if err != nil {
				return
			}
			defer conn.Close()
			conn.Write(data)
		}(peer)
	}

	return nil
}

// ============================================================================
// QUERY METHODS
// ============================================================================

// GetOnlinePeers returns all online peers
func (d *Discovery) GetOnlinePeers() []*models.Student {
	return d.peerRegistry.GetOnlinePeers()
}

// GetPeerCount returns the number of known peers
func (d *Discovery) GetPeerCount() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return len(d.knownPeers)
}

// IsPeerOnline checks if a peer is online
func (d *Discovery) IsPeerOnline(peerID string) bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	lastSeen, exists := d.knownPeers[peerID]
	if !exists {
		return false
	}

	return time.Since(lastSeen) < PeerTimeout
}

// GetStats returns discovery statistics
func (d *Discovery) GetStats() map[string]interface{} {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return map[string]interface{}{
		"is_running":   d.isRunning,
		"known_peers":  len(d.knownPeers),
		"online_peers": len(d.peerRegistry.GetOnlinePeers()),
		"subscribers":  len(d.subscribers),
	}
}
