package service

import (
	"context"
	"fmt"

	"myapp/internal/domain"
	"myapp/internal/repository"

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

const bcryptCost = 10

func (s *userService) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	if req.Name == "" || len(req.Name) < 2 {
		return nil, fmt.Errorf("имя должно быть минимум 2 символа")
	}
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("пароль должен быть минимум 6 символов")
	}

	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Avatar:   false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *userService) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}

func (s *userService) Update(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, fmt.Errorf("имя не может быть пустым")
		}
		user.Name = *req.Name
	}

	if req.Email != nil {
		if *req.Email == "" {
			return nil, fmt.Errorf("email не может быть пустым")
		}
		if *req.Email != user.Email {
			existing, _ := s.userRepo.GetByEmail(ctx, *req.Email)
			if existing != nil {
				return nil, domain.ErrEmailAlreadyExists
			}
		}
		user.Email = *req.Email
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	user.Password = ""
	return user, nil
}