package writer

import "context"

type Repository interface {
	Save(ctx context.Context, enterprise Enterprise) error
	Get(ctx context.Context, enterpriseKey EnterpriseKey) (EnterpriseData, error)
}
