package container

import (
	"github.com/IsaacDSC/search_content/internal/content/reader"
	"github.com/IsaacDSC/search_content/internal/content/writer"
)

type Handlers struct {
	WriterHandler writer.Handler
	ReaderHandler reader.Handler
}

func GetHandlers(services ServicesContainer) Handlers {
	wh := writer.NewHandler(services.WriterService)
	rh := reader.NewHandler(services.ReaderService)

	return Handlers{
		WriterHandler: wh,
		ReaderHandler: rh,
	}
}
