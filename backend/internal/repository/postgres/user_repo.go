package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (name, email, password, role, avatar, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(
		ctx, query,
		user.Name, user.Email, user.Password, user.Role, user.Avatar,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("ошибка создания: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	query := `SELECT id, name, email, role, avatar, last_seen_at, is_online, created_at FROM users ORDER BY id`

	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка: %w", err)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users SET
            name = COALESCE($1, name),
            email = COALESCE($2, email),
            role = COALESCE($3, role),
            avatar = COALESCE($4, avatar),
            updated_at = NOW()
        WHERE id = $5
        RETURNING updated_at
    `
	err := r.db.QueryRowContext(ctx, query,
		user.Name, user.Email, user.Role, user.Avatar, user.ID,
	).Scan(&user.UpdatedAt)

	if err == sql.ErrNoRows {
		return domain.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("ошибка обновления: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления пароля: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}