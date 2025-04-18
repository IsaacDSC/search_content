package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkFileWrite(b *testing.B) {
	dir, err := os.MkdirTemp("", "fs-benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	data := []byte("hello world")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := filepath.Join(dir, "test.txt")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			b.Fatal(err)
		}
		os.Remove(filename)
	}
}

func BenchmarkFileRead(b *testing.B) {
	dir, err := os.MkdirTemp("", "fs-benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filename := filepath.Join(dir, "test.txt")
	data := []byte("hello world")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := os.ReadFile(filename)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFileSeek(b *testing.B) {
	dir, err := os.MkdirTemp("", "fs-benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	filename := filepath.Join(dir, "test.txt")
	data := make([]byte, 1024*1024) // 1MB file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file, err := os.Open(filename)
		if err != nil {
			b.Fatal(err)
		}
		_, err = file.Seek(1024, 0)
		if err != nil {
			file.Close()
			b.Fatal(err)
		}
		file.Close()
	}
}

func BenchmarkDirectoryListing(b *testing.B) {
	dir, err := os.MkdirTemp("", "fs-benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create some files
	for i := 0; i < 10; i++ {
		filename := filepath.Join(dir, filepath.Base(filepath.Join("file", filepath.Clean(filepath.Join("file", "test.txt")))))
		if err := os.WriteFile(filename, []byte("test"), 0644); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		files, err := os.ReadDir(dir)
		if err != nil {
			b.Fatal(err)
		}
		_ = files
	}
}
