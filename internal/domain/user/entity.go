package user

import "time"

// User is the core domain entity, persisted as the "users" table.
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	Customer  string    `json:"customer" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// CreateUserRequest is the payload for creating a user.
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=6"`
	Customer string `json:"customer" validate:"required,min=2,max=100"`
}

// UpdateUserRequest is the payload for a partial update of a user.
// All fields are optional; only non-empty fields are applied.
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Password string `json:"password" validate:"omitempty,min=6"`
	Customer string `json:"customer" validate:"omitempty,min=2,max=100"`
}

// UserResponse is what gets returned to clients. Password is intentionally excluded.
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Customer  string    `json:"customer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse maps a domain User to its public response shape.
func ToResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Customer:  u.Customer,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
