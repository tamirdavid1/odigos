package fs

import (
	"os"
	"path/filepath"
)

func removeDir(hostDir string) error {
	return filepath.Walk(hostDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == hostDir {
			return nil
		}

		// Check if the file has a .so extension
		if filepath.Ext(info.Name()) == ".so" {
			return nil
		}

		// Remove the file or directory
		if info.IsDir() {
			return os.RemoveAll(path)
		} else {
			return os.Remove(path)
		}
	})
}
