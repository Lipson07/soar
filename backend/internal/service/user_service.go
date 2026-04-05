package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	exists, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if exists != nil {
		return nil, fmt.Errorf("ошибка почты")
	}

	exists, _ = s.userRepo.GetByUsername(ctx, req.Username)
	if exists != nil {
			return nil, fmt.Errorf("ошибка имени")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Status:    "offline",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) GetAll(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Status = status
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) Search(ctx context.Context, query string) ([]*domain.User, error) {
	return s.userRepo.Search(ctx, query)
}

func (s *userService) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.userRepo.Exists(ctx, id)
}