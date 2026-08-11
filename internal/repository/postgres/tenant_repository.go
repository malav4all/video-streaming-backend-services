package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/yourorg/go-user-service/internal/domain/tenant"
)

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) tenant.Repository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, t *tenant.Tenant) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *tenantRepository) GetByID(ctx context.Context, id string) (*tenant.Tenant, error) {
	var t tenant.Tenant
	if err := r.db.WithContext(ctx).First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tenantRepository) GetBySlug(ctx context.Context, slug string) (*tenant.Tenant, error) {
	var t tenant.Tenant
	if err := r.db.WithContext(ctx).First(&t, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tenantRepository) GetAll(ctx context.Context, limit, offset int) ([]tenant.Tenant, int64, error) {
	var tenants []tenant.Tenant
	var total int64

	if err := r.db.WithContext(ctx).Model(&tenant.Tenant{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Limit(limit).Offset(offset).
		Order("created_at desc").
		Find(&tenants).Error; err != nil {
		return nil, 0, err
	}
	return tenants, total, nil
}

func (r *tenantRepository) Update(ctx context.Context, t *tenant.Tenant) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&tenant.Tenant{}, "id = ?", id).Error
}