package writer

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestVideoInputDto_ToDomain(t *testing.T) {
	tests := []struct {
		name        string
		videoInput  VideoInputDto
		wantErr     bool
		wantDomain  Enterprise
		errContains string
	}{
		{
			name: "Valid video input",
			videoInput: VideoInputDto{
				VideoUrl:    "https://example.com/path/to/video.mp4",
				TambnailUrl: "https://example.com/path/to/thumbnail.jpg",
				Endpoint:    "/api/videos",
			},
			wantErr: false,
			wantDomain: func() Enterprise {
				u, _ := url.Parse("https://example.com/path/to/video.mp4")
				return Enterprise{
					Url:    u,
					Origin: "https://example.com",
					Paths:  []string{"path", "to", "video.mp4"},
					Path:   "/path/to/video.mp4",
					Video: Video{
						VideoUrl:    "https://example.com/path/to/video.mp4",
						TambnailUrl: "https://example.com/path/to/thumbnail.jpg",
					},
				}
			}(),
		},
		{
			name: "Video URL with query parameters",
			videoInput: VideoInputDto{
				VideoUrl:    "https://example.com/watch?v=abc123",
				TambnailUrl: "https://example.com/thumbs/abc123.jpg",
				Endpoint:    "/api/videos",
			},
			wantErr: false,
			wantDomain: func() Enterprise {
				u, _ := url.Parse("https://example.com/watch?v=abc123")
				return Enterprise{
					Url:    u,
					Origin: "https://example.com",
					Paths:  []string{"watch"},
					Path:   "/watch",
					Video: Video{
						VideoUrl:    "https://example.com/watch?v=abc123",
						TambnailUrl: "https://example.com/thumbs/abc123.jpg",
					},
				}
			}(),
		},
		{
			name: "Invalid video URL",
			videoInput: VideoInputDto{
				VideoUrl:    "http://invalid url with spaces",
				TambnailUrl: "https://example.com/thumb.jpg",
				Endpoint:    "/api/videos",
			},
			wantErr:     true,
			errContains: "parse",
		},
		{
			name: "Empty video URL",
			videoInput: VideoInputDto{
				VideoUrl:    "",
				TambnailUrl: "https://example.com/thumb.jpg",
				Endpoint:    "/api/videos",
			},
			wantErr:     true,
			errContains: "empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDomain, err := tt.videoInput.ToDomain()

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("VideoInputDto.ToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("VideoInputDto.ToDomain() error = %v, should contain %v", err, tt.errContains)
				return
			}

			// If we don't expect an error, validate the returned domain object
			if !tt.wantErr {
				// Compare URL strings since URL pointers won't be equal
				if gotDomain.Url.String() != tt.wantDomain.Url.String() {
					t.Errorf("VideoInputDto.ToDomain() Url = %v, want %v",
						gotDomain.Url.String(), tt.wantDomain.Url.String())
				}

				// Compare other fields
				if gotDomain.Origin != tt.wantDomain.Origin {
					t.Errorf("VideoInputDto.ToDomain() Origin = %v, want %v",
						gotDomain.Origin, tt.wantDomain.Origin)
				}

				if !reflect.DeepEqual(gotDomain.Paths, tt.wantDomain.Paths) {
					t.Errorf("VideoInputDto.ToDomain() Paths = %v, want %v",
						gotDomain.Paths, tt.wantDomain.Paths)
				}

				if gotDomain.Path != tt.wantDomain.Path {
					t.Errorf("VideoInputDto.ToDomain() Path = %v, want %v",
						gotDomain.Path, tt.wantDomain.Path)
				}

				if !reflect.DeepEqual(gotDomain.Video, tt.wantDomain.Video) {
					t.Errorf("VideoInputDto.ToDomain() Video = %v, want %v",
						gotDomain.Video, tt.wantDomain.Video)
				}
			}
		})
	}
}
