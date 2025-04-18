package repository

import (
	"context"
	"errors"
	"fmt"
	writer "github.com/IsaacDSC/search_content/internal/content/writer"
	"github.com/IsaacDSC/search_content/pkg/filesystem"
)

type FileSystemRepo struct {
	fsDrive filesystem.Driver
}

var _ writer.Repository = (*FileSystemRepo)(nil)

func NewFileSystemRepo(fsDrive filesystem.Driver) *FileSystemRepo {
	return &FileSystemRepo{fsDrive: fsDrive}
}

func (r FileSystemRepo) Save(ctx context.Context, enterprise writer.Enterprise) error {
	enterpriseKey := writer.NewEnterpriseKey(enterprise.Url)
	pathKey := writer.NewPathKey(enterprise.Url)

	result, err := r.Get(ctx, enterpriseKey)
	if !errors.Is(filesystem.ErrFileNotFound, err) && err != nil {
		return err
	}

	if errors.Is(filesystem.ErrFileNotFound, err) {
		data := writer.NewEnterprisesData(pathKey, enterprise)
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

func (r FileSystemRepo) Get(ctx context.Context, enterpriseKey writer.EnterpriseKey) (writer.EnterpriseData, error) {
	fileName := filesystem.NewFileName(enterpriseKey.String())
	data, err := r.fsDrive.Get(ctx, fileName)

	if err != nil {
		return writer.EnterpriseData{}, err
	}

	output, ok := data.(writer.EnterpriseData)
	if !ok {
		return writer.EnterpriseData{}, writer.ErrInvalidDataType
	}

	return output, nil
}
