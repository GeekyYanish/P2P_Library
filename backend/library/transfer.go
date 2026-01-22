/*
================================================================================
TRANSFER SERVICE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file handles file transfers between peers.

Go Concepts Used:
- Goroutines: Concurrent file transfers
- Channels: Progress reporting and control
- Interfaces: Transfer handler abstraction
- Error handling: Comprehensive error management
================================================================================
*/

package library

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"knowledge-exchange/utils"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// Transfer states
	TransferPending   = "pending"
	TransferActive    = "active"
	TransferCompleted = "completed"
	TransferFailed    = "failed"
	TransferCancelled = "cancelled"

	// Buffer size for file transfers
	TransferBufferSize = 32 * 1024 // 32 KB

	// Maximum concurrent transfers
	MaxConcurrentTransfers = 5
)

// ============================================================================
// TRANSFER REQUEST/RESPONSE STRUCTS
// ============================================================================

// TransferRequest represents a request to download a file
type TransferRequest struct {
	CID         string    `json:"cid"`
	RequesterID string    `json:"requester_id"`
	Timestamp   time.Time `json:"timestamp"`
}

// TransferResponse represents the response to a transfer request
type TransferResponse struct {
	CID      string `json:"cid"`
	Accepted bool   `json:"accepted"`
	Reason   string `json:"reason,omitempty"`
	FileSize int64  `json:"file_size"`
	Checksum string `json:"checksum"`
}

// Transfer represents an active file transfer
type Transfer struct {
	ID         string    `json:"id"`
	CID        string    `json:"cid"`
	FileName   string    `json:"file_name"`
	PeerID     string    `json:"peer_id"`
	Direction  string    `json:"direction"` // "upload" or "download"
	Status     string    `json:"status"`
	TotalBytes int64     `json:"total_bytes"`
	SentBytes  int64     `json:"sent_bytes"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time,omitempty"`
	Error      string    `json:"error,omitempty"`
	Progress   float64   `json:"progress"`
}

// ============================================================================
// PROGRESS REPORTING
// ============================================================================

// ProgressUpdate represents a transfer progress update
type ProgressUpdate struct {
	TransferID string
	BytesSent  int64
	TotalBytes int64
	Progress   float64
	Speed      float64 // bytes per second
}

// ============================================================================
// TRANSFER MANAGER
// ============================================================================

// TransferManager handles all file transfers
type TransferManager struct {
	// Active transfers
	transfers map[string]*Transfer

	// Progress channel for updates
	progressChan chan ProgressUpdate

	// Semaphore for concurrent transfer limiting
	semaphore chan struct{}

	// Mutex for thread-safe operations
	mutex sync.RWMutex

	// Indexer reference for file access
	indexer *Indexer

	// Stats
	totalUploads    int64
	totalDownloads  int64
	bytesUploaded   int64
	bytesDownloaded int64
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewTransferManager creates a new TransferManager
func NewTransferManager(indexer *Indexer) *TransferManager {
	return &TransferManager{
		transfers:    make(map[string]*Transfer),
		progressChan: make(chan ProgressUpdate, 100),
		semaphore:    make(chan struct{}, MaxConcurrentTransfers),
		indexer:      indexer,
	}
}

// ============================================================================
// UPLOAD METHODS
// ============================================================================

// HandleUploadRequest handles an incoming upload request
// Parameters:
//   - conn: The network connection
//   - request: The transfer request
//
// Returns:
//   - error: Error if upload fails
func (tm *TransferManager) HandleUploadRequest(conn net.Conn, request *TransferRequest) error {
	// Acquire semaphore (limits concurrent transfers)
	tm.semaphore <- struct{}{}
	defer func() { <-tm.semaphore }()

	// Get the file
	file, exists := tm.indexer.GetFile(request.CID)
	if !exists {
		return tm.sendResponse(conn, &TransferResponse{
			CID:      request.CID,
			Accepted: false,
			Reason:   "File not found",
		})
	}

	// Get local file path
	filePath, exists := tm.indexer.GetLocalFilePath(request.CID)
	if !exists {
		return tm.sendResponse(conn, &TransferResponse{
			CID:      request.CID,
			Accepted: false,
			Reason:   "File not available locally",
		})
	}

	// Create transfer record
	transfer := &Transfer{
		ID:         utils.HashString(fmt.Sprintf("%s-%d", request.CID, time.Now().UnixNano())),
		CID:        request.CID,
		FileName:   file.FileName,
		PeerID:     request.RequesterID,
		Direction:  "upload",
		Status:     TransferActive,
		TotalBytes: file.Size,
		StartTime:  time.Now(),
	}

	tm.addTransfer(transfer)
	defer tm.completeTransfer(transfer.ID)

	// Send acceptance response
	err := tm.sendResponse(conn, &TransferResponse{
		CID:      request.CID,
		Accepted: true,
		FileSize: file.Size,
		Checksum: file.Checksum,
	})
	if err != nil {
		transfer.Status = TransferFailed
		transfer.Error = err.Error()
		return err
	}

	// Open and stream file
	return tm.streamFile(conn, filePath, transfer)
}

// streamFile streams a file over the connection
func (tm *TransferManager) streamFile(conn net.Conn, filePath string, transfer *Transfer) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create buffer for reading
	buffer := make([]byte, TransferBufferSize)

	// Stream the file
	for {
		// Read from file
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			transfer.Status = TransferFailed
			transfer.Error = err.Error()
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Write to connection
		_, err = conn.Write(buffer[:n])
		if err != nil {
			transfer.Status = TransferFailed
			transfer.Error = err.Error()
			return fmt.Errorf("failed to send data: %w", err)
		}

		// Update progress
		transfer.SentBytes += int64(n)
		transfer.Progress = float64(transfer.SentBytes) / float64(transfer.TotalBytes) * 100

		// Send progress update
		tm.progressChan <- ProgressUpdate{
			TransferID: transfer.ID,
			BytesSent:  transfer.SentBytes,
			TotalBytes: transfer.TotalBytes,
			Progress:   transfer.Progress,
		}
	}

	transfer.Status = TransferCompleted
	transfer.EndTime = time.Now()
	tm.totalUploads++
	tm.bytesUploaded += transfer.TotalBytes

	return nil
}

// ============================================================================
// DOWNLOAD METHODS
// ============================================================================

// Download downloads a file from a remote peer
// Parameters:
//   - peerAddress: The address of the peer (ip:port)
//   - cid: The Content Identifier of the file
//   - savePath: Where to save the downloaded file
//   - requesterID: ID of the requesting peer
//
// Returns:
//   - error: Error if download fails
func (tm *TransferManager) Download(peerAddress, cid, savePath, requesterID string) error {
	// Acquire semaphore
	tm.semaphore <- struct{}{}
	defer func() { <-tm.semaphore }()

	// Connect to peer
	conn, err := utils.Connect(peerAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to peer: %w", err)
	}
	defer conn.Close()

	// Send transfer request
	request := &TransferRequest{
		CID:         cid,
		RequesterID: requesterID,
		Timestamp:   time.Now(),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	msg := &utils.Message{
		Type:    utils.MsgTypeRequest,
		Sender:  requesterID,
		Payload: requestData,
	}

	if err := utils.SendMessage(conn, msg); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Receive response
	responseMsg, err := utils.ReceiveMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}

	var response TransferResponse
	if err := json.Unmarshal(responseMsg.Payload, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Accepted {
		return fmt.Errorf("transfer rejected: %s", response.Reason)
	}

	// Create transfer record
	transfer := &Transfer{
		ID:         utils.HashString(fmt.Sprintf("%s-%d", cid, time.Now().UnixNano())),
		CID:        cid,
		PeerID:     peerAddress,
		Direction:  "download",
		Status:     TransferActive,
		TotalBytes: response.FileSize,
		StartTime:  time.Now(),
	}

	tm.addTransfer(transfer)
	defer tm.completeTransfer(transfer.ID)

	// Receive file
	return tm.receiveFile(conn, savePath, response.FileSize, response.Checksum, transfer)
}

// receiveFile receives a file from a connection
func (tm *TransferManager) receiveFile(conn net.Conn, savePath string, fileSize int64, checksum string, transfer *Transfer) error {
	// Create output file
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Buffer for receiving
	buffer := make([]byte, TransferBufferSize)
	var received int64

	// Receive data
	for received < fileSize {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			transfer.Status = TransferFailed
			transfer.Error = err.Error()
			return fmt.Errorf("failed to receive data: %w", err)
		}

		// Write to file
		_, err = file.Write(buffer[:n])
		if err != nil {
			transfer.Status = TransferFailed
			transfer.Error = err.Error()
			return fmt.Errorf("failed to write data: %w", err)
		}

		received += int64(n)
		transfer.SentBytes = received
		transfer.Progress = float64(received) / float64(fileSize) * 100

		// Send progress update
		tm.progressChan <- ProgressUpdate{
			TransferID: transfer.ID,
			BytesSent:  received,
			TotalBytes: fileSize,
			Progress:   transfer.Progress,
		}
	}

	// Verify checksum
	computedChecksum, err := utils.HashFile(savePath)
	if err != nil {
		return fmt.Errorf("failed to compute checksum: %w", err)
	}

	if computedChecksum != checksum {
		// Remove corrupted file
		os.Remove(savePath)
		transfer.Status = TransferFailed
		transfer.Error = "Checksum verification failed"
		return fmt.Errorf("checksum verification failed")
	}

	transfer.Status = TransferCompleted
	transfer.EndTime = time.Now()
	tm.totalDownloads++
	tm.bytesDownloaded += fileSize

	return nil
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// sendResponse sends a transfer response
func (tm *TransferManager) sendResponse(conn net.Conn, response *TransferResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	msg := &utils.Message{
		Type:    utils.MsgTypeResponse,
		Sender:  "local",
		Payload: data,
	}

	return utils.SendMessage(conn, msg)
}

// addTransfer adds a transfer to the active list
func (tm *TransferManager) addTransfer(t *Transfer) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.transfers[t.ID] = t
}

// completeTransfer marks a transfer as complete
func (tm *TransferManager) completeTransfer(id string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if t, exists := tm.transfers[id]; exists {
		if t.Status == TransferActive {
			t.Status = TransferCompleted
		}
		t.EndTime = time.Now()
	}
}

// GetTransfer returns a transfer by ID
func (tm *TransferManager) GetTransfer(id string) (*Transfer, bool) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	t, exists := tm.transfers[id]
	return t, exists
}

// GetActiveTransfers returns all active transfers
func (tm *TransferManager) GetActiveTransfers() []*Transfer {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	var active []*Transfer
	for _, t := range tm.transfers {
		if t.Status == TransferActive {
			active = append(active, t)
		}
	}
	return active
}

// GetProgressChannel returns the progress update channel
func (tm *TransferManager) GetProgressChannel() <-chan ProgressUpdate {
	return tm.progressChan
}

// GetStats returns transfer statistics
func (tm *TransferManager) GetStats() map[string]interface{} {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	return map[string]interface{}{
		"total_uploads":    tm.totalUploads,
		"total_downloads":  tm.totalDownloads,
		"bytes_uploaded":   tm.bytesUploaded,
		"bytes_downloaded": tm.bytesDownloaded,
		"active_transfers": len(tm.GetActiveTransfers()),
	}
}

// CancelTransfer cancels an active transfer
func (tm *TransferManager) CancelTransfer(id string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	t, exists := tm.transfers[id]
	if !exists {
		return fmt.Errorf("transfer not found: %s", id)
	}

	if t.Status != TransferActive {
		return fmt.Errorf("transfer is not active")
	}

	t.Status = TransferCancelled
	t.EndTime = time.Now()
	return nil
}
