package writer

import (
	"errors"
	"github.com/IsaacDSC/search_content/internal/content/entity"
	"net/url"
	"strings"
)

type VideoInputDto struct {
	VideoUrl    string `json:"video_url"`
	TambnailUrl string `json:"thumbnail_url"`
	Endpoint    string `json:"endpoint"`
}

func (v *VideoInputDto) ToDomain() (entity.Enterprise, error) {
	endpoint, err := url.Parse(v.VideoUrl)
	if err != nil {
		return entity.Enterprise{}, err
	}

	origin := endpoint.Scheme + "://" + endpoint.Host
	path := endpoint.Path
	paths := strings.Split(path, "/")[1:] // Remove a primeira barra

	if v.VideoUrl == "" {
		return entity.Enterprise{}, errors.New("video url is empty")
	}

	if v.TambnailUrl == "" {
		return entity.Enterprise{}, errors.New("thumbnail url is empty")
	}

	return entity.Enterprise{
		Url:    endpoint,
		Origin: origin,
		Paths:  paths,
		Path:   path,
		Video: entity.Video{
			VideoUrl:    v.VideoUrl,
			TambnailUrl: v.TambnailUrl,
		},
	}, nil
}
