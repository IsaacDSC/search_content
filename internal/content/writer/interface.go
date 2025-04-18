package writer

import "context"

type Repository interface {
	Save(ctx context.Context, enterprise Enterprise) error
}
