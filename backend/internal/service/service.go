package service

import (
	"context"
	"myapp/internal/domain"
)

type UserService interface {
	Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
}