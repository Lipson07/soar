package repository

import (
	"context"
	"myapp/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error
	Delete(ctx context.Context, id int64) error
}
type ChatRepository interface {
	Create(ctx context.Context, user *domain.Chat) error
	GetByID(ctx context.Context, id int64) (*domain.Chat, error)
	GetByName(ctx context.Context, email string) (*domain.Chat, error)
	GetAll(ctx context.Context) ([]domain.Chat, error)
	Update(ctx context.Context, user *domain.Chat) error
	Delete(ctx context.Context, id int64) error
}