package fs

import (
	"os"
	"path/filepath"

	"github.com/odigos-io/odigos/odiglet/pkg/log"
)

func removeFilesInDir(hostDir string) error {
	return filepath.Walk(hostDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Logger.Info("Error accessing path %s: %v\n", path, err)
			return nil // Continue walking
		}

		// Skip the root directory itself
		if path == hostDir {
			return nil
		}

		// Check if the entry is a directory
		if info.IsDir() {
			log.Logger.Info("Skipping directory: %s\n", path)
			return nil
		}

		// Check if the file has a .so extension
		if filepath.Ext(info.Name()) == ".so" {
			log.Logger.Info("Skipping .so file: %s\n", path)
			return nil
		}

		// Remove the file
		log.Logger.Info("Removing file: %s\n", path)
		err = os.Remove(path)
		if err != nil {
			log.Logger.Error(err, "Error removing file")
		}

		return nil
	})
}
