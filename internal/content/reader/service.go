package reader

import (
	"context"
	"fmt"
	"github.com/IsaacDSC/search_content/internal/content/entity"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s Service) GetContent(ctx context.Context, endpoint EndpointDto) (entity.Video, error) {
	url, err := endpoint.ToDomain()
	if err != nil {
		return entity.Video{}, err
	}

	key := entity.NewEnterpriseKey(url)
	fmt.Println(key)

	//data, err := s.repository.Get(ctx, key)
	//if err != nil {
	//	return entity.Video{}, err
	//}

	return entity.Video{}, nil
}
