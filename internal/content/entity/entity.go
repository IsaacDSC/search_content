package entity

import (
	"net/url"
	"strings"
)

type Video struct {
	VideoUrl    string
	TambnailUrl string
}

func (v Video) IsEmpty() bool {
	var empty Video
	return empty == v
}

type Enterprise struct {
	Url    *url.URL
	Origin string
	Paths  []string
	Path   string
	Video  Video
}

type EnterpriseKey string

func (ek EnterpriseKey) String() string {
	return string(ek)
}

func NewEnterpriseKey(u *url.URL) EnterpriseKey {
	return EnterpriseKey(u.Scheme + "://" + u.Host)
}

type PathKey string

func NewPathKey(u *url.URL) PathKey {
	//pathSanitized := strings.ReplaceAll(u.Path, "*", "")
	return PathKey(strings.ReplaceAll(u.Path, "//", "/"))
}

func (pk PathKey) ToListPaths() (output []string) {
	r := strings.Split(string(pk), "/")
	for i := range r {
		if r[i] == "" || r[i] == "/" {
			continue
		}

		output = append(output, r[i])
	}

	return
}
