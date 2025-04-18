package repository

import (
	"context"
	"errors"
	"github.com/IsaacDSC/search_content/internal/content/writer"
	"github.com/IsaacDSC/search_content/pkg/filesystem"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/url"
	"testing"
)

func TestFileSystemRepo_Save(t *testing.T) {
	// Helper function to parse URLs in test cases
	parseURL := func(rawURL string) *url.URL {
		u, err := url.Parse(rawURL)
		if err != nil {
			t.Fatalf("Failed to parse URL %q: %v", rawURL, err)
		}
		return u
	}

	// Define test cases
	tests := []struct {
		name          string
		setupMock     func(*gomock.Controller) (filesystem.Driver, writer.Enterprise)
		expectedError error
	}{
		{
			name: "save new enterprise",
			setupMock: func(ctrl *gomock.Controller) (filesystem.Driver, writer.Enterprise) {
				mockDriver := filesystem.NewMockDriver(ctrl)
				enterprise := writer.Enterprise{
					Url: parseURL("https://example.com/video"),
				}
				enterpriseKey := writer.NewEnterpriseKey(enterprise.Url)
				pathKey := writer.NewPathKey(enterprise.Url)
				fileName := filesystem.NewFileName(enterpriseKey.String())

				mockDriver.EXPECT().
					Get(gomock.Any(), fileName).
					Return(nil, filesystem.ErrFileNotFound)

				expectedData := writer.NewEnterprisesData(pathKey, enterprise)
				mockDriver.EXPECT().
					Save(gomock.Any(), fileName, expectedData).
					Return(nil)

				return mockDriver, enterprise
			},
			expectedError: nil,
		},
		{
			name: "update existing enterprise",
			setupMock: func(ctrl *gomock.Controller) (filesystem.Driver, writer.Enterprise) {
				mockDriver := filesystem.NewMockDriver(ctrl)
				enterprise := writer.Enterprise{
					Url: parseURL("https://example.com/video"),
				}
				enterpriseKey := writer.NewEnterpriseKey(enterprise.Url)
				pathKey := writer.NewPathKey(enterprise.Url)
				fileName := filesystem.NewFileName(enterpriseKey.String())

				existingData := writer.NewEnterprisesData(pathKey, enterprise)

				mockDriver.EXPECT().
					Get(gomock.Any(), fileName).
					Return(existingData, nil)

				mockDriver.EXPECT().
					Save(gomock.Any(), fileName, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ filesystem.FileName, data any) error {
						// Verify the data being saved has the enterprise
						savedData, ok := data.(writer.EnterpriseData)
						assert.True(t, ok, "data should be EnterpriseData")

						savedEnterprise, exists := savedData[pathKey]
						assert.True(t, exists, "enterprise should exist in data")
						assert.Equal(t, enterprise.Url.String(), savedEnterprise.Url.String())

						return nil
					})

				return mockDriver, enterprise
			},
			expectedError: nil,
		},
		{
			name: "handle driver error",
			setupMock: func(ctrl *gomock.Controller) (filesystem.Driver, writer.Enterprise) {
				mockDriver := filesystem.NewMockDriver(ctrl)
				enterprise := writer.Enterprise{
					Url: parseURL("https://example.com/video"),
				}
				enterpriseKey := writer.NewEnterpriseKey(enterprise.Url)
				fileName := filesystem.NewFileName(enterpriseKey.String())
				expectedErr := errors.New("driver error")

				mockDriver.EXPECT().
					Get(gomock.Any(), fileName).
					Return(nil, expectedErr)

				return mockDriver, enterprise
			},
			expectedError: errors.New("driver error"),
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDriver, enterprise := tt.setupMock(ctrl)
			repo := NewFileSystemRepo(mockDriver)

			// Execute
			err := repo.Save(context.Background(), enterprise)

			// Assert
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if tt.expectedError.Error() != "" {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			}
		})
	}
}

func TestFileSystemRepo_Get(t *testing.T) {
	// Helper function to parse URLs in test cases
	parseURL := func(rawURL string) *url.URL {
		u, err := url.Parse(rawURL)
		if err != nil {
			t.Fatalf("Failed to parse URL %q: %v", rawURL, err)
		}
		return u
	}

	// Define test cases
	tests := []struct {
		name           string
		setupMock      func(*gomock.Controller) (filesystem.Driver, writer.EnterpriseKey)
		expectedData   writer.EnterpriseData
		expectedError  error
		errorValidator func(t *testing.T, err error)
	}{
		{
			name: "get existing enterprise data",
			setupMock: func(ctrl *gomock.Controller) (filesystem.Driver, writer.EnterpriseKey) {
				mockDriver := filesystem.NewMockDriver(ctrl)
				url := parseURL("https://example.com/video")
				enterpriseKey := writer.NewEnterpriseKey(url)
				pathKey := writer.NewPathKey(url)
				fileName := filesystem.NewFileName(enterpriseKey.String())

				enterprise := writer.Enterprise{
					Url: url,
				}
				expectedData := writer.NewEnterprisesData(pathKey, enterprise)

				mockDriver.EXPECT().
					Get(gomock.Any(), fileName).
					Return(expectedData, nil)

				return mockDriver, enterpriseKey
			},
			expectedData:  writer.NewEnterprisesData(writer.NewPathKey(parseURL("https://example.com/video")), writer.Enterprise{Url: parseURL("https://example.com/video")}),
			expectedError: nil,
		},
		{
			name: "handle file not found",
			setupMock: func(ctrl *gomock.Controller) (filesystem.Driver, writer.EnterpriseKey) {
				mockDriver := filesystem.NewMockDriver(ctrl)
				url := parseURL("https://example.com/video")
				enterpriseKey := writer.NewEnterpriseKey(url)
				fileName := filesystem.NewFileName(enterpriseKey.String())

				mockDriver.EXPECT().
					Get(gomock.Any(), fileName).
					Return(nil, filesystem.ErrFileNotFound)

				return mockDriver, enterpriseKey
			},
			expectedData:  writer.EnterpriseData{},
			expectedError: filesystem.ErrFileNotFound,
			errorValidator: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, filesystem.ErrFileNotFound)
			},
		},
		{
			name: "handle invalid data type",
			setupMock: func(ctrl *gomock.Controller) (filesystem.Driver, writer.EnterpriseKey) {
				mockDriver := filesystem.NewMockDriver(ctrl)
				url := parseURL("https://example.com/video")
				enterpriseKey := writer.NewEnterpriseKey(url)
				fileName := filesystem.NewFileName(enterpriseKey.String())

				// Return string instead of EnterpriseData to trigger type assertion failure
				invalidData := "not an enterprise data"

				mockDriver.EXPECT().
					Get(gomock.Any(), fileName).
					Return(invalidData, nil)

				return mockDriver, enterpriseKey
			},
			expectedData:  writer.EnterpriseData{},
			expectedError: writer.ErrInvalidDataType,
			errorValidator: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, writer.ErrInvalidDataType)
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDriver, enterpriseKey := tt.setupMock(ctrl)
			repo := NewFileSystemRepo(mockDriver)

			// Execute
			data, err := repo.Get(context.Background(), enterpriseKey)

			// Assert
			if tt.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, data)
			} else {
				assert.Error(t, err)
				assert.Empty(t, data)
				if tt.errorValidator != nil {
					tt.errorValidator(t, err)
				}
			}
		})
	}
}
