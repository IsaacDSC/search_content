package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileSystem struct {
	mu      sync.RWMutex
	baseDir string
}

func NewFileSystem() *FileSystem {
	// Get the current working directory as base directory
	cwd, err := os.Getwd()
	if err != nil {
		// Fallback to relative path if we can't get working directory
		cwd = "."
	}
	return &FileSystem{
		baseDir: cwd,
	}
}

type FileName string

func NewFileName(input string) FileName {
	return FileName(fmt.Sprintf("assets/tmp/%s_json", input))
}

func (fn FileName) String() string {
	return string(fn)
}

// getFullPath returns the absolute path for a given filename
func (fs *FileSystem) getFullPath(filename FileName) string {
	return filepath.Join(fs.baseDir, filename.String())
}

// ensureDir ensures that the directory for the given file path exists
func (fs *FileSystem) ensureDir(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// FileExists checks if a file already exists
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
