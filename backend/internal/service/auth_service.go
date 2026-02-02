package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/shahnajsc/OnePointLedger/backend/internal/model"
	"github.com/shahnajsc/OnePointLedger/backend/internal/repo"
)

type AuthService struct {
	users     *repo.UserRepo
	jwtSecret []byte
}

func NewAuthService(users *repo.UserRepo, jwtSecret string) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	return s.users.CreateUser(ctx, email, string(hash))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub": u.ID,
		"email": u.Email,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
