package fs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/odigos-io/odigos/odiglet/pkg/log"
)

const (
	// we don't want to overload the CPU so we limit to small number of goroutines
	maxWorkers = 4

	// 32 KB buffer for I/O operations
	bufferSize = 32 * 1024
)

// getNumberOfWorkers returns the number of workers to use for copying files.
// It returns the minimum of maxWorkers and the number of CPUs divided by 4.
func getNumberOfWorkers() int {
	return min(maxWorkers, max(1, runtime.NumCPU()/4))
}

func copyDirectories(srcDir string, destDir string) error {
	start := time.Now()
	hasOSFiles := HasSOFiles(destDir)

	files, err := getFiles(srcDir, hasOSFiles)
	if err != nil {
		return err
	}

	// Create the destination directory if it doesn't exist
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Create a buffered channel to control concurrency
	numWorkers := getNumberOfWorkers()
	fileChan := make(chan string, numWorkers)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(fileChan, srcDir, destDir, &wg)
	}

	// Send files to the channel
	for _, file := range files {
		fileChan <- file
	}

	// Close the channel and wait for workers to finish
	close(fileChan)
	wg.Wait()
	log.Logger.V(0).Info("Finished copying instrumentation files to host", "duration", time.Since(start))
	return nil
}

func worker(fileChan <-chan string, sourceDir, destDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Allocate a buffer once per goroutine.
	buf := make([]byte, bufferSize)
	for file := range fileChan {
		destFile := filepath.Join(destDir, file[len(sourceDir)+1:])
		err := copyFile(file, destFile, buf)
		if err != nil {
			log.Logger.Error(err, "error copying file", "file", file)
		}
	}
}

func getFiles(dir string, hasOSFiles bool) ([]string, error) {

	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if hasOSFiles {
				switch ext := filepath.Ext(path); ext {
				case ".so", ".node", "node.d", ".a":
					log.Logger.Info(fmt.Sprintf("Skipping .so files : %s", path))
					return nil
				}
			}
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func copyFile(src, dst string, buf []byte) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file and directories if needed
	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy using the provided buffer.
	for {
		n, err := srcFile.Read(buf)
		if n > 0 {
			if _, err := dstFile.Write(buf[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func HasSOFiles(dir string) bool {
	var hasSOFiles bool
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		switch ext := filepath.Ext(info.Name()); ext {
		case ".so", ".node", "node.d", ".a", ".o":
			hasSOFiles = true
			return filepath.SkipDir
		}
		return nil
	})

	return hasSOFiles
}
