package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/odigos-io/odigos/odiglet/pkg/log"
)

func removeFilesInDir(hostDir string) error {
	return filepath.Walk(hostDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Logger.Info("Error accessing path  %s: %v\n", path, err)
			return nil // Continue walking
		}

		// Skip the root directory itself]
		if path == hostDir {
			return nil
		}

		// Check if the entry is a directory
		if info.IsDir() {
			log.Logger.Info(fmt.Sprintf("Skipping directory: %s", path))
			return nil
		}

		// Check if the file has a .so extension
		switch ext := filepath.Ext(info.Name()); ext {
		case ".so", ".node", "node.d", ".a":
			log.Logger.Info(fmt.Sprintf("Skipping .so files : %s", path))
			return nil
		}

		// Remove the file
		log.Logger.Info(fmt.Sprintf("Removing file: %s", path))
		err = os.Remove(path)
		if err != nil {
			log.Logger.Error(err, "Error removing file")
		}

		return nil
	})
}
