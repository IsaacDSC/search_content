package writer

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_Register(t *testing.T) {
	tests := []struct {
		name        string
		input       VideoInputDto
		setupMocks  func(mockRepo *MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful registration",
			input: VideoInputDto{
				VideoUrl:    "https://example.com/video.mp4",
				TambnailUrl: "https://example.com/thumbnail.jpg",
				Endpoint:    "/api/videos",
			},
			setupMocks: func(mockRepo *MockRepository) {
				mockRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid input",
			input: VideoInputDto{
				VideoUrl:    "",
				TambnailUrl: "https://example.com/thumbnail.jpg",
				Endpoint:    "/api/videos",
			},
			setupMocks: func(mockRepo *MockRepository) {
				// No repository calls expected for invalid input
			},
			wantErr:     true,
			errContains: "empty",
		},
		{
			name: "repository save error",
			input: VideoInputDto{
				VideoUrl:    "https://example.com/video.mp4",
				TambnailUrl: "https://example.com/thumbnail.jpg",
				Endpoint:    "/api/videos",
			},
			setupMocks: func(mockRepo *MockRepository) {
				mockRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr:     true,
			errContains: "failed to save entity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create the mock repository
			mockRepo := NewMockRepository(ctrl)

			// Setup expectations based on the test case
			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			// Create service with mock repository
			service := NewContentUseCase(mockRepo)

			// Execute the method being tested
			err := service.Register(context.Background(), tt.input)

			// Check if error expectations match using assert
			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got nil")
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains),
						"Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
			}
		})
	}
}
