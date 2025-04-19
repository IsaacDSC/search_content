package reader_test

import (
	"github.com/IsaacDSC/search_content/internal/content/builder"
	"testing"

	"github.com/IsaacDSC/search_content/internal/content/entity"
	"github.com/IsaacDSC/search_content/internal/content/reader"
	"github.com/stretchr/testify/assert"
)

func TestEnterpriseData_GetContent(t *testing.T) {
	videoBuilder := builder.NewVideoBuilder()

	videoA := videoBuilder.WithRandomData().Build()
	videoB := videoBuilder.WithRandomData().Build()
	videoC := videoBuilder.WithRandomData().Build()

	enterpriseBuilder := builder.NewEnterpriseBuilder()

	tests := []struct {
		name           string
		enterpriseData reader.EnterpriseData
		input          entity.PathKey
		expectedVideo  entity.Video
		expectedFound  bool
	}{
		{
			name: "Exact match",
			enterpriseData: reader.EnterpriseData{
				entity.PathKey("path/to/video"): enterpriseBuilder.WithRandomData().
					WithPath("path/to/video").
					WithVideo(videoA).
					Build(),
			},
			input:         entity.PathKey("path/to/video"),
			expectedVideo: videoA,
			expectedFound: true,
		},
		{
			name: "No match",
			enterpriseData: reader.EnterpriseData{
				entity.PathKey("path/to/video"): enterpriseBuilder.WithRandomData().
					WithPath("path/to/video").
					WithVideo(videoA).
					Build(),
			},
			input:         entity.PathKey("different/path"),
			expectedVideo: entity.Video{},
			expectedFound: false,
		},
		{
			name: "Wildcard match",
			enterpriseData: reader.EnterpriseData{
				entity.PathKey("path/*/video"): enterpriseBuilder.WithRandomData().
					WithPath("path/*/video").
					WithVideo(videoB).
					Build(),
			},
			input:         entity.PathKey("path/something/video"),
			expectedVideo: videoB,
			expectedFound: true,
		},
		{
			name: "Multiple entries with one match",
			enterpriseData: reader.EnterpriseData{
				entity.PathKey("path/to/videoA"): enterpriseBuilder.WithRandomData().
					WithPath("path/to/videoA").
					WithVideo(videoA).
					Build(),
				entity.PathKey("path/to/videoB"): enterpriseBuilder.WithRandomData().
					WithPath("path/to/videoB").
					WithVideo(videoB).
					Build(),
			},
			input:         entity.PathKey("path/to/videoB"),
			expectedVideo: videoB,
			expectedFound: true,
		},
		{
			name: "Multiple wildcard patterns",
			enterpriseData: reader.EnterpriseData{
				entity.PathKey("path/*/videoA"): enterpriseBuilder.WithRandomData().
					WithPath("path/*/videoA").
					WithVideo(videoA).
					Build(),
				entity.PathKey("path/to/*"): enterpriseBuilder.WithRandomData().
					WithPath("path/to/*").
					WithVideo(videoC).
					Build(),
			},
			input:         entity.PathKey("path/to/something"),
			expectedVideo: videoC,
			expectedFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			video, found := tt.enterpriseData.GetContent(tt.input)

			assert.Equal(t, tt.expectedFound, found)
			if tt.expectedFound {
				assert.Equal(t, tt.expectedVideo, video)
			}
		})
	}
}
