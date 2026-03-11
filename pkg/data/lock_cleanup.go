package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// CleanupStaleLocks removes stale lock files from the data directory
// This is needed when the process was killed forcefully (SIGKILL) and didn't
// get a chance to clean up properly.
//
// Lock files are created by Lucene/Diagon IndexWriter to prevent concurrent writes.
// They should be automatically removed on clean shutdown, but may remain if:
// - Process was killed with SIGKILL (-9)
// - Process crashed
// - System crashed
//
// This function should be called during startup BEFORE loading shards.
func CleanupStaleLocks(dataDir string, logger *zap.Logger) error {
	logger.Info("Checking for stale lock files", zap.String("data_dir", dataDir))

	// Check if data directory exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		logger.Debug("Data directory does not exist, no locks to clean")
		return nil
	}

	lockFilesRemoved := 0

	// Walk the entire data directory tree looking for *.lock files
	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Look for .lock files (Lucene/Diagon write locks)
		if strings.HasSuffix(info.Name(), ".lock") {
			logger.Warn("Found stale lock file, removing",
				zap.String("path", path),
				zap.String("file", info.Name()))

			if err := os.Remove(path); err != nil {
				logger.Error("Failed to remove lock file",
					zap.String("path", path),
					zap.Error(err))
				return fmt.Errorf("failed to remove lock file %s: %w", path, err)
			}

			lockFilesRemoved++
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to clean up stale locks: %w", err)
	}

	if lockFilesRemoved > 0 {
		logger.Info("Removed stale lock files",
			zap.Int("count", lockFilesRemoved))
	} else {
		logger.Debug("No stale lock files found")
	}

	return nil
}
