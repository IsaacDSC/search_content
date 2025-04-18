// Package filesystem provides a simple file-based storage system
// for JSON data with concurrency safety.
package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	once sync.Once
	fs   *FileSystem
)

// FileSystem provides thread-safe file operations for storing
// and retrieving data using the local filesystem.
type FileSystem struct {
	mu      sync.RWMutex // protects concurrent access to files
	baseDir string       // base directory for all file operations
}

// NewFileSystem creates a singleton instance of FileSystem.
// It uses the current working directory as the base directory.
// If getting the working directory fails, it falls back to ".".
func NewFileSystem() *FileSystem {
	once.Do(func() {
		// Get the current working directory as base directory
		cwd, err := os.Getwd()
		if err != nil {
			// Fallback to relative path if we can't get working directory
			cwd = "."
		}

		fs = &FileSystem{
			baseDir: cwd,
		}
	})

	return fs
}

// FileName is a type for file names used in the filesystem.
// It encapsulates the naming convention for files.
type FileName string

// NewFileName creates a new FileName with the specified input
// and applies standard filepath formatting.
func NewFileName(input string) FileName {
	return FileName(fmt.Sprintf("assets/tmp/%s_json", input))
}

// String returns the string representation of the FileName.
func (fn FileName) String() string {
	return string(fn)
}

// getFullPath returns the absolute path for a given filename.
func (fs *FileSystem) getFullPath(filename FileName) string {
	return filepath.Join(fs.baseDir, filename.String())
}

// ensureDir ensures that the directory for the given file path exists.
// It creates the directory if it doesn't exist yet.
func (fs *FileSystem) ensureDir(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// FileExists checks if a file already exists.
// Returns true if the file exists, false otherwise.
// Any errors other than "file not exists" are returned.
func (fs *FileSystem) FileExists(ctx context.Context, key FileName) (bool, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	fullPath := fs.getFullPath(key)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Save stores data to a file specified by key.
// The data can be any type and will be marshaled to JSON
// unless it's already a string, in which case it's stored directly.
// It ensures the target directory exists before writing.
func (fs *FileSystem) Save(ctx context.Context, key FileName, data any) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fullPath := fs.getFullPath(key)

	// Ensure target directory exists
	if err := fs.ensureDir(fullPath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	var bytes []byte
	var err error

	// Check if data is already a string
	if str, ok := data.(string); ok {
		bytes = []byte(str)
	} else {
		// Marshal the data to JSON if it's not already a string
		bytes, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
	}

	// Use context to potentially handle timeout
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := os.WriteFile(fullPath, bytes, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

// Get retrieves data from a file specified by key.
// It reads the file and unmarshals the JSON content into a map[string]any.
// The context can be used for cancellation or timeout.
func (fs *FileSystem) Get(ctx context.Context, key FileName) (any, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	fullPath := fs.getFullPath(key)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		b, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		output := make(map[string]any)
		if err := json.Unmarshal(b, &output); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file: %w", err)
		}

		return output, nil
	}
}
