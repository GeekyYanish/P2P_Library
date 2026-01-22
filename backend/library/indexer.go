/*
================================================================================
INDEXER SERVICE - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file handles file indexing and Content Identifier (CID) generation.

Go Concepts Used:
- Goroutines: Concurrent file scanning
- Channels: Communication between goroutines
- Maps: File index storage
- Slices: Dynamic arrays for file lists
- Error handling: Go's error patterns
================================================================================
*/

package library

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"knowledge-exchange/models"
	"knowledge-exchange/utils"
)

// ============================================================================
// INDEXER SERVICE STRUCT
// ============================================================================

// Indexer manages the file index for the P2P network
type Indexer struct {
	// fileIndex stores all indexed files by CID
	fileIndex *models.FileIndex

	// localFiles stores files available on this peer
	localFiles map[string]string // CID -> file path

	// watchDir is the directory being watched for new files
	watchDir string

	// mutex for thread-safe operations
	mutex sync.RWMutex

	// isRunning indicates if the indexer is active
	isRunning bool

	// stopChan signals the watcher to stop
	stopChan chan struct{}
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewIndexer creates a new Indexer instance
// Parameters:
//   - watchDir: Directory to watch for shared files
func NewIndexer(watchDir string) *Indexer {
	return &Indexer{
		fileIndex:  models.NewFileIndex(),
		localFiles: make(map[string]string),
		watchDir:   watchDir,
		isRunning:  false,
		stopChan:   make(chan struct{}),
	}
}

// ============================================================================
// INDEXING METHODS
// ============================================================================

// IndexFile indexes a single file and generates its CID
// Parameters:
//   - filePath: Path to the file to index
//   - ownerID: ID of the file owner
//
// Returns:
//   - *models.AcademicFile: The indexed file info
//   - error: Error if indexing fails
func (idx *Indexer) IndexFile(filePath, ownerID string) (*models.AcademicFile, error) {
	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Check file size
	if fileInfo.Size() > utils.MaxFileSizeBytes {
		return nil, fmt.Errorf("file exceeds maximum size limit")
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if !utils.IsAllowedFileType(ext) {
		return nil, fmt.Errorf("file type %s is not allowed", ext)
	}

	// Read file content for CID generation
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create the academic file record
	academicFile := models.NewAcademicFile(
		fileInfo.Name(),
		ownerID,
		fileInfo.Size(),
		ext,
		content,
	)

	// Add to index
	idx.mutex.Lock()
	idx.fileIndex.Add(academicFile)
	idx.localFiles[academicFile.CID] = filePath
	idx.mutex.Unlock()

	return academicFile, nil
}

// ScanDirectory scans a directory and indexes all valid files
// Uses Goroutines for concurrent processing
// Parameters:
//   - dirPath: Directory to scan
//   - ownerID: ID of the file owner
//
// Returns:
//   - []*models.AcademicFile: List of indexed files
//   - error: Error if scan fails
func (idx *Indexer) ScanDirectory(dirPath, ownerID string) ([]*models.AcademicFile, error) {
	var indexedFiles []*models.AcademicFile
	var mu sync.Mutex

	// Create a channel for file paths
	fileChan := make(chan string, 100)
	resultChan := make(chan *models.AcademicFile, 100)
	errorChan := make(chan error, 100)

	// Create a WaitGroup to wait for all goroutines
	var wg sync.WaitGroup

	// Worker goroutine count
	workerCount := 4

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				file, err := idx.IndexFile(filePath, ownerID)
				if err != nil {
					errorChan <- err
					continue
				}
				resultChan <- file
			}
		}()
	}

	// Goroutine to walk directory and send file paths
	go func() {
		filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if !info.IsDir() {
				fileChan <- path
			}
			return nil
		})
		close(fileChan)
	}()

	// Goroutine to collect results
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results
	for file := range resultChan {
		mu.Lock()
		indexedFiles = append(indexedFiles, file)
		mu.Unlock()
	}

	return indexedFiles, nil
}

// ============================================================================
// LOOKUP METHODS
// ============================================================================

// GetFile retrieves a file by its CID
func (idx *Indexer) GetFile(cid string) (*models.AcademicFile, bool) {
	return idx.fileIndex.Get(cid)
}

// GetLocalFilePath returns the local path for a CID
func (idx *Indexer) GetLocalFilePath(cid string) (string, bool) {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	path, exists := idx.localFiles[cid]
	return path, exists
}

// Search searches for files matching a query
func (idx *Indexer) Search(query string) []*models.AcademicFile {
	return idx.fileIndex.Search(query)
}

// GetBySubject returns files for a specific subject
func (idx *Indexer) GetBySubject(subject string) []*models.AcademicFile {
	return idx.fileIndex.GetBySubject(subject)
}

// GetAllFiles returns all indexed files
func (idx *Indexer) GetAllFiles() []*models.AcademicFile {
	return idx.fileIndex.GetAllFiles()
}

// GetLocalFiles returns all locally available files
func (idx *Indexer) GetLocalFiles() []*models.AcademicFile {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	var files []*models.AcademicFile
	for cid := range idx.localFiles {
		if file, exists := idx.fileIndex.Get(cid); exists {
			files = append(files, file)
		}
	}
	return files
}

// ============================================================================
// FILE OPERATIONS
// ============================================================================

// GetFileContent reads and returns the content of a local file
func (idx *Indexer) GetFileContent(cid string) ([]byte, error) {
	path, exists := idx.GetLocalFilePath(cid)
	if !exists {
		return nil, fmt.Errorf("file not found locally: %s", cid)
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}

// RemoveFile removes a file from the index
func (idx *Indexer) RemoveFile(cid string) error {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	// Remove from both indexes
	idx.fileIndex.Remove(cid)
	delete(idx.localFiles, cid)

	return nil
}

// ============================================================================
// DIRECTORY WATCHER
// ============================================================================

// StartWatcher starts watching the shared files directory for changes
// Uses Goroutines for background monitoring
func (idx *Indexer) StartWatcher(ownerID string, interval time.Duration) {
	if idx.isRunning {
		return
	}

	idx.isRunning = true

	// Start watcher goroutine
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Rescan directory for new files
				idx.ScanDirectory(idx.watchDir, ownerID)
			case <-idx.stopChan:
				return
			}
		}
	}()
}

// StopWatcher stops the directory watcher
func (idx *Indexer) StopWatcher() {
	if idx.isRunning {
		idx.stopChan <- struct{}{}
		idx.isRunning = false
	}
}

// ============================================================================
// STATISTICS
// ============================================================================

// GetStats returns indexer statistics
func (idx *Indexer) GetStats() map[string]interface{} {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	return map[string]interface{}{
		"total_indexed": idx.fileIndex.Count(),
		"local_files":   len(idx.localFiles),
		"is_watching":   idx.isRunning,
		"watch_dir":     idx.watchDir,
	}
}
