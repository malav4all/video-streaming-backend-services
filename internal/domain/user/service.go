package user

import "context"

// Service is the "port" the HTTP layer (or any other delivery mechanism,
// e.g. gRPC, CLI) depends on. Concrete business logic lives in the usecase layer.
type Service interface {
	Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	GetByID(ctx context.Context, id uint) (*UserResponse, error)
	GetAll(ctx context.Context, page, pageSize int) ([]UserResponse, int64, error)
	Update(ctx context.Context, id uint, req *UpdateUserRequest) (*UserResponse, error)
	Delete(ctx context.Context, id uint) error
}
