// Driver defines the interface for file-based storage operations.
// This interface helps with testing and dependency inversion.
package filesystem

import (
	"context"
)

// Driver defines the operations available for interacting with
// a file-based storage system.
type Driver interface {
	// FileExists checks if a file exists in the filesystem.
	// Returns true if the file exists, false otherwise.
	FileExists(ctx context.Context, key FileName) (bool, error)

	// Save stores data to a file specified by key.
	// The data can be any type and will be marshaled to JSON.
	Save(ctx context.Context, key FileName, data any) error

	// Get retrieves data from a file specified by key.
	// It reads the file and unmarshals the JSON content.
	Get(ctx context.Context, key FileName) (any, error)
}

// Ensure FileSystem implements the Driver interface
var _ Driver = (*FileSystem)(nil)
