package tenant

import "context"

type Repository interface {
	Create(ctx context.Context, t *Tenant) error
	GetByID(ctx context.Context, id string) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	GetAll(ctx context.Context, limit, offset int) ([]Tenant, int64, error)
	Update(ctx context.Context, t *Tenant) error
	Delete(ctx context.Context, id string) error
}
