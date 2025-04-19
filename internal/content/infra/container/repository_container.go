package container

import (
	"github.com/IsaacDSC/search_content/internal/content/infra/repository"
	"github.com/IsaacDSC/search_content/pkg/filesystem"
)

type RepositoryContainer struct {
	Repository repository.Repository
}

func NewRepositoryContainer() RepositoryContainer {
	fsDriver := filesystem.NewFileSystem()
	repo := repository.NewFileSystemRepo(fsDriver)
	return RepositoryContainer{
		Repository: repo,
	}
}
