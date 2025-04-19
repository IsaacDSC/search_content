package builder

import (
	"github.com/IsaacDSC/search_content/internal/content/entity"
	"github.com/go-faker/faker/v4"
)

// VideoBuilder provides a fluent interface to construct a Video entity
type VideoBuilder struct {
	video entity.Video
}

// NewVideoBuilder creates a new VideoBuilder instance
func NewVideoBuilder() *VideoBuilder {
	video := entity.Video{
		VideoUrl:    faker.URL(),
		TambnailUrl: faker.URL(),
	}

	return &VideoBuilder{
		video: video,
	}
}

func (vb *VideoBuilder) WithRandomData() *VideoBuilder {
	vb.video.VideoUrl = faker.URL()
	vb.video.TambnailUrl = faker.URL()
	return vb
}

// WithVideoUrl sets the VideoUrl field
func (vb *VideoBuilder) WithVideoUrl(videoUrl string) *VideoBuilder {
	vb.video.VideoUrl = videoUrl
	return vb
}

// WithThumbnailUrl sets the TambnailUrl field
// Note: Following the field name from entity definition
func (vb *VideoBuilder) WithThumbnailUrl(thumbnailUrl string) *VideoBuilder {
	vb.video.TambnailUrl = thumbnailUrl
	return vb
}

// Build returns the constructed Video entity
func (vb *VideoBuilder) Build() entity.Video {
	return vb.video
}
