package builder

import (
	"github.com/IsaacDSC/search_content/internal/content/entity"
	"github.com/go-faker/faker/v4"
	"net/url"
	"strings"
)

// EnterpriseBuilder provides a fluent interface to construct an Enterprise entity
type EnterpriseBuilder struct {
	enterprise entity.Enterprise
}

// NewEnterpriseBuilder creates a new EnterpriseBuilder instance
func NewEnterpriseBuilder() *EnterpriseBuilder {
	randomURL, _ := url.Parse(faker.URL())

	enterprise := entity.Enterprise{
		Url:    randomURL,
		Origin: randomURL.Scheme + "://" + randomURL.Host,
		Path:   randomURL.Path,
		Paths:  strings.Split(strings.TrimPrefix(randomURL.Path, "/"), "/"),
		Video:  entity.Video{},
	}

	return &EnterpriseBuilder{
		enterprise: enterprise,
	}
}

// WithRandomData populates the Enterprise with random data
func (eb *EnterpriseBuilder) WithRandomData() *EnterpriseBuilder {
	randomURL, _ := url.Parse(faker.URL())

	eb.enterprise.Url = randomURL
	eb.enterprise.Origin = randomURL.Scheme + "://" + randomURL.Host
	eb.enterprise.Path = randomURL.Path
	eb.enterprise.Paths = strings.Split(strings.TrimPrefix(randomURL.Path, "/"), "/")
	eb.enterprise.Video = NewVideoBuilder().WithRandomData().Build()

	return eb
}

// WithUrl sets the Url field
func (eb *EnterpriseBuilder) WithUrl(url *url.URL) *EnterpriseBuilder {
	eb.enterprise.Url = url
	return eb
}

// WithOrigin sets the Origin field
func (eb *EnterpriseBuilder) WithOrigin(origin string) *EnterpriseBuilder {
	eb.enterprise.Origin = origin
	return eb
}

// WithPaths sets the Paths field
func (eb *EnterpriseBuilder) WithPaths(paths []string) *EnterpriseBuilder {
	eb.enterprise.Paths = paths
	return eb
}

// WithPath sets the Path field
func (eb *EnterpriseBuilder) WithPath(path string) *EnterpriseBuilder {
	eb.enterprise.Path = path
	eb.enterprise.Url.Path = path
	return eb
}

// WithVideo sets the Video field
func (eb *EnterpriseBuilder) WithVideo(video entity.Video) *EnterpriseBuilder {
	eb.enterprise.Video = video
	return eb
}

// Build returns the constructed Enterprise entity
func (eb *EnterpriseBuilder) Build() entity.Enterprise {
	return eb.enterprise
}
