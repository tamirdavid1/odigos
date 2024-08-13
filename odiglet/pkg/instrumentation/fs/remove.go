package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/odigos-io/odigos/odiglet/pkg/log"
)

func removeFilesInDir(hostDir string, filesToKeep map[string]struct{}) error {
	shouldRecreateCFiles := ShouldRecreateAllCFiles()
	log.Logger.V(0).Info(fmt.Sprintf("Removing files in directory: %s, shouldRecreateCFiles: %t", hostDir, shouldRecreateCFiles))

	// Mark directories as protected if they contain a file that needs to be preserved.
	// If C files should be recreated, skip marking any directories as protected.
	protectedDirs := make(map[string]bool)
	if !shouldRecreateCFiles {
		for file := range filesToKeep {
			dir := filepath.Dir(file)
			for dir != hostDir {
				protectedDirs[dir] = true
				dir = filepath.Dir(dir)
			}
			protectedDirs[hostDir] = true // Protect the main directory
		}
	}

	return filepath.Walk(hostDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if path == hostDir {
			return nil
		}

		// Skip any files listed in filesToKeepMap
		if !info.IsDir() {
			if _, found := filesToKeep[path]; found {
				log.Logger.V(0).Info(fmt.Sprintf("Skipping protected file: %s", path))
				return nil
			}
		}

		// Skip protected directories
		if info.IsDir() && protectedDirs[path] {
			log.Logger.V(0).Info(fmt.Sprintf("Skipping protected directory: %s", path))
			return nil
		}

		// Remove unprotected files and directories
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("error removing %s: %w", path, err)
		}

		// Skip further processing in this directory since it has been removed
		if info.IsDir() {
			return filepath.SkipDir
		}

		return nil
	})
}
