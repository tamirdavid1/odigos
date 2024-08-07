package fs

import (
	"os"
	"path/filepath"

	"github.com/odigos-io/odigos/odiglet/pkg/log"
)

func removeFilesInDir(hostDir string) error {
	return filepath.Walk(hostDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself]
		if path == hostDir {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !ShouldRecreateAllCFiles() {
			log.Logger.Info("Removing the file: %s, because Should ShouldRecreateAllCFiles is %s", path, ShouldRecreateAllCFiles())
			switch ext := filepath.Ext(info.Name()); ext {
			case ".so", ".node", "node.d", ".a":
				log.Logger.Info("Skipping removing the file: %s", path)
				return nil
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
