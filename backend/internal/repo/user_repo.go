package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shahnajsc/OnePointLedger/backend/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, email, passwordHash string) (model.User, error) {
	const q = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id::text, email, password_hash, created_at;
	`

	var u model.User
	err := r.db.QueryRowContext(ctx, q, email, passwordHash).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (model.User, error) {
	const q = `
		SELECT id::text, email, password_hash, created_at
		FROM users
		WHERE email = $1;
	`

	var u model.User
	err := r.db.QueryRowContext(ctx, q, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, sql.ErrNoRows
	}
	return u, err
}
