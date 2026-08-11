package tenant

import "context"

type Service interface {
	Create(ctx context.Context, req *CreateTenantRequest) (*TenantResponse, error)
	GetByID(ctx context.Context, id string) (*TenantResponse, error)
	GetBySlug(ctx context.Context, slug string) (*TenantResponse, error)
	GetAll(ctx context.Context, page, pageSize int) ([]TenantResponse, int64, error)
	Update(ctx context.Context, id string, req *UpdateTenantRequest) (*TenantResponse, error)
	Delete(ctx context.Context, id string) error
}