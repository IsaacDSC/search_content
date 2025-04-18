package writer

import (
	"net/url"
	"strings"
)

type Video struct {
	VideoUrl    string
	TambnailUrl string
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
	pathSanitized := strings.ReplaceAll(u.Path, "*", "")
	return PathKey(strings.ReplaceAll(pathSanitized, "//", "/"))
}
