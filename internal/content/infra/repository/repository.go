package repository

import (
	"github.com/IsaacDSC/search_content/internal/content/reader"
	"github.com/IsaacDSC/search_content/internal/content/writer"
)

type Repository interface {
	reader.Repository
	writer.Repository
}
