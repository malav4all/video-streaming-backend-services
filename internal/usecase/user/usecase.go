package usecase

import (
	"context"
	"errors"

	"github.com/yourorg/go-user-service/internal/domain/user"
	"github.com/yourorg/go-user-service/pkg/hash"
	"gorm.io/gorm"
)

// ErrUserNotFound is returned when a requested user does not exist.
var ErrUserNotFound = errors.New("user not found")

type userUsecase struct {
	repo user.Repository
}

// NewUserUsecase wires a repository implementation into the user.Service port.
func NewUserUsecase(repo user.Repository) user.Service {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Create(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	hashed, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		Name:     req.Name,
		Password: hashed,
		Customer: req.Customer,
	}

	if err := u.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return user.ToResponse(newUser), nil
}

func (u *userUsecase) GetByID(ctx context.Context, id uint) (*user.UserResponse, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user.ToResponse(existing), nil
}

func (u *userUsecase) GetAll(ctx context.Context, page, pageSize int) ([]user.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	users, total, err := u.repo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	res := make([]user.UserResponse, 0, len(users))
	for i := range users {
		res = append(res, *user.ToResponse(&users[i]))
	}

	return res, total, nil
}

func (u *userUsecase) Update(ctx context.Context, id uint, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Customer != "" {
		existing.Customer = req.Customer
	}
	if req.Password != "" {
		hashed, err := hash.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		existing.Password = hashed
	}

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return user.ToResponse(existing), nil
}

func (u *userUsecase) Delete(ctx context.Context, id uint) error {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return u.repo.Delete(ctx, id)
}
