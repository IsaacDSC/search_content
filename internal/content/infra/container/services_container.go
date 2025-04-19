package container

import (
	"github.com/IsaacDSC/search_content/internal/content/reader"
	"github.com/IsaacDSC/search_content/internal/content/writer"
)

type ServicesContainer struct {
	WriterService writer.Service
	ReaderService reader.Service
}

func NewServicesContainer(repositories RepositoryContainer) ServicesContainer {
	writerService := writer.NewContentUseCase(repositories.Repository)
	readerService := reader.NewContentUseCase(repositories.Repository)

	return ServicesContainer{
		WriterService: writerService,
		ReaderService: readerService,
	}
}
