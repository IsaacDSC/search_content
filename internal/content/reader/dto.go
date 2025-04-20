package reader

import (
	"net/url"
)

type EndpointDto string

func NewEndpointDto(endpoint string) EndpointDto {
	return EndpointDto(endpoint)
}

func (e EndpointDto) ToDomain() (*url.URL, error) {
	endpoint, err := url.Parse(string(e))
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}
