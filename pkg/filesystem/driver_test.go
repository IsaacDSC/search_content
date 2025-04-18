package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestFileSystem_WriteAndRead(t *testing.T) {
	// Setup temporary directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	tmpDir := filepath.Join(cwd, "assets", "tmp")
	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize filesystem
	fs := NewFileSystem()
	ctx := context.Background()

	// Data to save
	testData := `{"name":"test","value":123}`
	fileName := NewFileName("test")

	// Write data
	err = fs.Save(ctx, fileName, testData)
	if err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}

	// Read data back
	data, err := fs.Get(ctx, fileName)
	if err != nil {
		t.Fatalf("Failed to get data: %v", err)
	}

	// Verify data
	result, ok := data.(map[string]interface{})
	if !ok {
		t.Fatal("Invalid data type returned")
	}

	if result["name"] != "test" || result["value"].(float64) != 123 {
		t.Fatalf("Data mismatch. Got: %v", result)
	}
}

func TestFileSystem_Concurrency(t *testing.T) {
	// Setup temporary directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	tmpDir := filepath.Join(cwd, "assets", "tmp")
	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize filesystem
	fs := NewFileSystem()
	ctx := context.Background()

	// Number of concurrent operations
	const numOps = 10
	var wg sync.WaitGroup
	wg.Add(numOps * 2) // Write + Read operations

	// Lock for coordinating verification
	var mu sync.Mutex
	results := make(map[string]bool)

	// Perform concurrent writes and reads
	for i := 0; i < numOps; i++ {
		go func(index int) {
			defer wg.Done()

			// Create unique data for each goroutine
			data := map[string]interface{}{
				"id":    index,
				"value": index * 10,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				t.Errorf("Failed to marshal data: %v", err)
				return
			}

			fileName := NewFileName(fmt.Sprintf("concurrent_%d", index))

			// Save data
			err = fs.Save(ctx, fileName, string(jsonData))
			if err != nil {
				t.Errorf("Failed to save data: %v", err)
				return
			}
		}(i)

		go func(index int) {
			defer wg.Done()

			// Wait briefly to allow write operation to complete
			time.Sleep(10 * time.Millisecond)

			fileName := NewFileName(fmt.Sprintf("concurrent_%d", index))

			// Read data
			data, err := fs.Get(ctx, fileName)
			if err != nil {
				t.Errorf("Failed to read data: %v", err)
				return
			}

			result, ok := data.(map[string]interface{})
			if !ok {
				t.Error("Invalid data type returned")
				return
			}

			// Verify data
			id, idOk := result["id"].(float64)
			value, valueOk := result["value"].(float64)

			if !idOk || !valueOk || int(id) != index || int(value) != index*10 {
				t.Errorf("Data mismatch for index %d. Got: %v", index, result)
				return
			}

			mu.Lock()
			results[fmt.Sprintf("concurrent_%d", index)] = true
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	// Verify all operations completed successfully
	if len(results) != numOps {
		t.Fatalf("Expected %d successful operations, got %d", numOps, len(results))
	}
}
