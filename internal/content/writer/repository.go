package writer

import (
	"context"
	"github.com/IsaacDSC/search_content/internal/content/entity"
)

type Repository interface {
	Save(ctx context.Context, enterprise entity.Enterprise) error
}
