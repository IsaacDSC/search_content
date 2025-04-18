package writer

import (
	"context"
	"fmt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Register(ctx context.Context, input VideoInputDto) error {
	entity, err := input.ToDomain()
	if err != nil {
		return err
	}

	if err = s.repository.Save(ctx, entity); err != nil {
		return fmt.Errorf("failed to save entity: %w", err)
	}

	return nil
}
