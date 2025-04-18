package writer

import (
	"context"
	"errors"
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

	err = s.repository.Save(ctx, entity)
	if !errors.Is(err, ErrAlreadyRegistered) && err != nil {
		return err
	}

	return nil
}
