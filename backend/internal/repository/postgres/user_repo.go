package postgres

import (
	"backend/internal/domain"
	"context"
	"database/sql"

	"github.com/google/uuid"
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
		INSERT INTO users (id, username, email, password_hash, avatar_url, status, last_seen, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.Password,
		user.AvatarURL, user.Status, user.LastSeen, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, status, last_seen, created_at, updated_at FROM users WHERE id = $1`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, status, last_seen, created_at, updated_at FROM users WHERE email = $1`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, status, last_seen, created_at, updated_at FROM users WHERE username = $1`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, status, last_seen, created_at, updated_at FROM users ORDER BY created_at DESC`
	var users []*domain.User
	err := r.db.SelectContext(ctx, &users, query)
	return users, err
}

func (r *UserRepository) Search(ctx context.Context, query string) ([]*domain.User, error) {
	searchQuery := `
		SELECT id, username, email, password_hash, avatar_url, status, last_seen, created_at, updated_at 
		FROM users 
		WHERE username ILIKE $1 OR email ILIKE $1 
		ORDER BY username 
		LIMIT 20
	`
	var users []*domain.User
	err := r.db.SelectContext(ctx, &users, searchQuery, "%"+query+"%")
	return users, err
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET username = $1, email = $2, password_hash = $3, 
		avatar_url = $4, status = $5, last_seen = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.Password,
		user.AvatarURL, user.Status, user.LastSeen, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)
	return exists, err
}