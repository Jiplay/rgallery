package scanner

import (
	"sync"
)

var (
	// ScanInProgress tracks whether a scan is currently in progress
	scanInProgress bool
	// mutex to protect access to scanInProgress
	scanMutex sync.RWMutex
)

// IsScanInProgress returns true if a scan is currently in progress
func IsScanInProgress() bool {
	scanMutex.RLock()
	defer scanMutex.RUnlock()
	return scanInProgress
}

// SetScanInProgress sets the scan in progress status
func SetScanInProgress(status bool) {
	scanMutex.Lock()
	defer scanMutex.Unlock()
	scanInProgress = status
}
