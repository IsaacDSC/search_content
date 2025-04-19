package writer

import (
	"context"
	"fmt"
)

type Service interface {
	Register(ctx context.Context, input VideoInputDto) error
}

type ContentUseCase struct {
	repository Repository
}

func NewContentUseCase(repository Repository) *ContentUseCase {
	return &ContentUseCase{repository: repository}
}

func (s *ContentUseCase) Register(ctx context.Context, input VideoInputDto) error {
	entity, err := input.ToDomain()
	if err != nil {
		return err
	}

	if err = s.repository.Save(ctx, entity); err != nil {
		return fmt.Errorf("failed to save entity: %w", err)
	}

	return nil
}
