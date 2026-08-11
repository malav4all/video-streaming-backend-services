package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/yourorg/go-user-service/internal/domain/tenant"
)

var ErrTenantNotFound = errors.New("tenant not found")

type tenantUsecase struct {
	repo tenant.Repository
}

func NewTenantUsecase(repo tenant.Repository) tenant.Service {
	return &tenantUsecase{repo: repo}
}

func (u *tenantUsecase) Create(ctx context.Context, req *tenant.CreateTenantRequest) (*tenant.TenantResponse, error) {
	newTenant := &tenant.Tenant{
		ID:   uuid.NewString(),
		Slug: req.Slug,
		Name: req.Name,
	}
	if err := u.repo.Create(ctx, newTenant); err != nil {
		return nil, err
	}
	return tenant.ToResponse(newTenant), nil
}

func (u *tenantUsecase) GetByID(ctx context.Context, id string) (*tenant.TenantResponse, error) {
	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return tenant.ToResponse(t), nil
}

func (u *tenantUsecase) GetBySlug(ctx context.Context, slug string) (*tenant.TenantResponse, error) {
	t, err := u.repo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return tenant.ToResponse(t), nil
}

func (u *tenantUsecase) GetAll(ctx context.Context, page, pageSize int) ([]tenant.TenantResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	tenants, total, err := u.repo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	res := make([]tenant.TenantResponse, 0, len(tenants))
	for i := range tenants {
		res = append(res, *tenant.ToResponse(&tenants[i]))
	}
	return res, total, nil
}

func (u *tenantUsecase) Update(ctx context.Context, id string, req *tenant.UpdateTenantRequest) (*tenant.TenantResponse, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		existing.Name = req.Name
	}

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return tenant.ToResponse(existing), nil
}

func (u *tenantUsecase) Delete(ctx context.Context, id string) error {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTenantNotFound
		}
		return err
	}
	return u.repo.Delete(ctx, id)
}
