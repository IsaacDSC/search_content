package writer

import (
	"errors"
	"net/url"
	"strings"
)

type VideoInputDto struct {
	VideoUrl    string `json:"video_url"`
	TambnailUrl string `json:"thumbnail_url"`
	Endpoint    string `json:"endpoint"`
}

func (v *VideoInputDto) ToDomain() (Enterprise, error) {
	endpoint, err := url.Parse(v.VideoUrl)
	if err != nil {
		return Enterprise{}, err
	}

	origin := endpoint.Scheme + "://" + endpoint.Host
	path := endpoint.Path
	paths := strings.Split(path, "/")[1:] // Remove a primeira barra

	if v.VideoUrl == "" {
		return Enterprise{}, errors.New("video url is empty")
	}

	if v.TambnailUrl == "" {
		return Enterprise{}, errors.New("thumbnail url is empty")
	}

	return Enterprise{
		Url:    endpoint,
		Origin: origin,
		Paths:  paths,
		Path:   path,
		Video: Video{
			VideoUrl:    v.VideoUrl,
			TambnailUrl: v.TambnailUrl,
		},
	}, nil
}
