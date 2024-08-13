package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/odigos-io/odigos/odiglet/pkg/log"
)

func removeFilesInDir(hostDir string) error {
	shouldRecreateCFiles := ShouldRecreateAllCFiles()
	log.Logger.V(0).Info(fmt.Sprintf("Removing files in directory: %s, shouldRecreateCFiles: %s", hostDir, fmt.Sprintf("%t", shouldRecreateCFiles)))

	return filepath.Walk(hostDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself]
		if path == hostDir {
			return nil
		}

		// Skip directories containing "ebpf" in their name
		if info.IsDir() && strings.Contains(info.Name(), "ebpf") {
			return nil
		}

		// Remove and skip other directories
		if info.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("error removing directory %s: %w", path, err)
			}
			return filepath.SkipDir
		}

		if !shouldRecreateCFiles {
			// filter out C files in ebpf directories
			if strings.Contains(filepath.Dir(path), "ebpf") {
				switch ext := filepath.Ext(info.Name()); ext {
				case ".so", ".node", "node.d", ".a":
					return nil
				}
			}
		}

		// Remove the file
		err = os.Remove(path)
		if err != nil {
			return err
		}

		return nil
	})
}
