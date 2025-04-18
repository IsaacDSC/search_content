package writer

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNewEnterpriseKey(t *testing.T) {
	tests := []struct {
		name     string
		urlInput string
		want     EnterpriseKey
	}{
		{
			name:     "Simple URL",
			urlInput: "https://example.com/path",
			want:     EnterpriseKey("https://example.com"),
		},
		{
			name:     "URL with port",
			urlInput: "http://localhost:8080/api/v1",
			want:     EnterpriseKey("http://localhost:8080"),
		},
		{
			name:     "URL with username and password",
			urlInput: "https://user:pass@example.com/secure",
			want:     EnterpriseKey("https://example.com"),
		},
		{
			name:     "URL with query parameters",
			urlInput: "https://api.example.com/search?q=test",
			want:     EnterpriseKey("https://api.example.com"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.urlInput)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}
			if got := NewEnterpriseKey(u); got != tt.want {
				t.Errorf("NewEnterpriseKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPathKey(t *testing.T) {
	tests := []struct {
		name     string
		urlInput string
		want     PathKey
	}{
		{
			name:     "Simple path",
			urlInput: "https://example.com/path/to/resource",
			want:     PathKey("/path/to/resource"),
		},
		{
			name:     "Path with wildcard",
			urlInput: "https://example.com/api/v1/users/*",
			want:     PathKey("/api/v1/users/"),
		},
		{
			name:     "Root path",
			urlInput: "https://example.com/",
			want:     PathKey("/"),
		},
		{
			name:     "Path with multiple wildcards",
			urlInput: "https://example.com/*/products/*",
			want:     PathKey("/products/"),
		},
		{
			name:     "Path with query parameters",
			urlInput: "https://example.com/search?q=test&page=1",
			want:     PathKey("/search"),
		},
		{
			name:     "Complex path with multiple wildcards",
			urlInput: "https://example.com/api/*/users/*/profile",
			want:     PathKey("/api/users/profile"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.urlInput)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}
			got := NewPathKey(u)
			if got != tt.want {
				t.Errorf("NewPathKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnterpriseStructure(t *testing.T) {
	// Test the Enterprise struct by constructing an instance and checking its fields
	urlStr := "https://example.com/path/to/resource"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}

	enterprise := Enterprise{
		Url:    u,
		Origin: "https://example.com",
		Paths:  []string{"path", "to", "resource"},
		Path:   "/path/to/resource",
		Video: Video{
			VideoUrl:    "https://example.com/path/to/resource/video.mp4",
			TambnailUrl: "https://example.com/path/to/resource/thumbnail.jpg",
		},
	}

	// Verify fields using subtests for better organization
	t.Run("URL field", func(t *testing.T) {
		if enterprise.Url.String() != urlStr {
			t.Errorf("Enterprise.Url = %v, want %v", enterprise.Url.String(), urlStr)
		}
	})

	t.Run("Origin field", func(t *testing.T) {
		if enterprise.Origin != "https://example.com" {
			t.Errorf("Enterprise.Origin = %v, want %v", enterprise.Origin, "https://example.com")
		}
	})

	t.Run("Paths field", func(t *testing.T) {
		expectedPaths := []string{"path", "to", "resource"}
		if !reflect.DeepEqual(enterprise.Paths, expectedPaths) {
			t.Errorf("Enterprise.Paths = %v, want %v", enterprise.Paths, expectedPaths)
		}
	})

	t.Run("Path field", func(t *testing.T) {
		if enterprise.Path != "/path/to/resource" {
			t.Errorf("Enterprise.Path = %v, want %v", enterprise.Path, "/path/to/resource")
		}
	})

	t.Run("Video field", func(t *testing.T) {
		expectedVideo := Video{
			VideoUrl:    "https://example.com/path/to/resource/video.mp4",
			TambnailUrl: "https://example.com/path/to/resource/thumbnail.jpg",
		}
		if !reflect.DeepEqual(enterprise.Video, expectedVideo) {
			t.Errorf("Enterprise.Video = %v, want %v", enterprise.Video, expectedVideo)
		}
	})
}
