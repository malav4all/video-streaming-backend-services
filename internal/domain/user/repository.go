package user

import "context"

// Repository is the "port" that any persistence adapter (Postgres, Mongo, etc.)
// must implement. The domain and usecase layers depend only on this interface,
// never on a concrete database driver.
type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetAll(ctx context.Context, limit, offset int) ([]User, int64, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uint) error
}
