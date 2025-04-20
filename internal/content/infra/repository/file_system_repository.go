package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IsaacDSC/search_content/internal/content/entity"
	"github.com/IsaacDSC/search_content/internal/content/reader"
	"github.com/IsaacDSC/search_content/internal/content/writer"
	"github.com/IsaacDSC/search_content/pkg/filesystem"
)

type FileSystemRepo struct {
	fsDrive filesystem.Driver
}

var _ Repository = (*FileSystemRepo)(nil)

func NewFileSystemRepo(fsDrive filesystem.Driver) *FileSystemRepo {
	return &FileSystemRepo{fsDrive: fsDrive}
}

func (r FileSystemRepo) Save(ctx context.Context, enterprise entity.Enterprise) error {
	enterpriseKey := entity.NewEnterpriseKey(enterprise.Url)
	pathKey := entity.NewPathKey(enterprise.Url)

	result, err := r.Get(ctx, enterpriseKey)
	if !errors.Is(filesystem.ErrFileNotFound, err) && err != nil {
		return err
	}

	if errors.Is(filesystem.ErrFileNotFound, err) {
		data := reader.NewEnterprisesData(pathKey, enterprise)
		fileName := filesystem.NewFileName(enterpriseKey.String())
		return r.fsDrive.Save(ctx, fileName, data)
	}

	fileName := filesystem.NewFileName(enterpriseKey.String())
	data := result.Append(pathKey, enterprise)
	if err := r.fsDrive.Save(ctx, fileName, data); err != nil {
		return fmt.Errorf("failed to save enterprise file: %w", err)
	}

	return nil
}

func (r FileSystemRepo) Get(ctx context.Context, enterpriseKey entity.EnterpriseKey) (reader.EnterpriseData, error) {
	fileName := filesystem.NewFileName(enterpriseKey.String())
	data, err := r.fsDrive.Get(ctx, fileName)

	if err != nil {
		return reader.EnterpriseData{}, err
	}

	output, ok := data.(map[string]any)
	if !ok {
		return reader.EnterpriseData{}, writer.ErrInvalidDataType
	}

	// Convert map to JSON bytes
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		return reader.EnterpriseData{}, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Unmarshal JSON into EnterpriseData
	var enterpriseData reader.EnterpriseData
	if err := json.Unmarshal(jsonBytes, &enterpriseData); err != nil {
		return reader.EnterpriseData{}, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return enterpriseData, nil
}
