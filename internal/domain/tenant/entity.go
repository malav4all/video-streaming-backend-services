package tenant

import "time"

// Tenant is the core domain entity, persisted as the "tenants" table.
// Slug matches the corresponding Keycloak group name (e.g. "tenant-acme"),
// which is how the middleware resolves a token's groups claim to a real tenant.
type Tenant struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey"`
	Slug      string    `json:"slug" gorm:"type:varchar(100);uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"type:varchar(150);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Tenant) TableName() string {
	return "tenants"
}

type CreateTenantRequest struct {
	Slug string `json:"slug" validate:"required,min=2,max=100"`
	Name string `json:"name" validate:"required,min=2,max=150"`
}

type UpdateTenantRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=150"`
}

type TenantResponse struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToResponse(t *Tenant) *TenantResponse {
	return &TenantResponse{
		ID:        t.ID,
		Slug:      t.Slug,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
