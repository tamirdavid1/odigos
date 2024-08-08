package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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
			log.Logger.Info(fmt.Sprintf("removing the file: %s, because Should ShouldRecreateAllCFiles is %s", path, strconv.FormatBool(ShouldRecreateAllCFiles())))
			switch ext := filepath.Ext(info.Name()); ext {
			case ".so", ".node", "node.d", ".a":
				log.Logger.Info(fmt.Sprintf("Skipping removing the file: %s", path))
				return nil
			}
		}

		return nil
	})
}
