package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"context"
	"errors"

	"github.com/shahnajsc/OnePointLedger/backend/internal/service"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	auth *service.AuthService
}

func NewHandler(auth *service.AuthService) *Handler {
	return &Handler{auth: auth}
}

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c creds) valid() bool {
	return strings.Contains(c.Email, "@") && len(c.Password) >= 8
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := contextWithTimeout(r, 5*time.Second)
	defer cancel()

	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if !c.valid() {
		http.Error(w, "invalid email or password (min 8 chars)", http.StatusBadRequest)
		return
	}

	u, err := h.auth.Signup(ctx, c.Email, c.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, "email already registered", http.StatusConflict)
			return
		}
		http.Error(w, "could not create user", http.StatusBadRequest)
		return
	}

	resp := map[string]any{
		"id":        u.ID,
		"email":     u.Email,
		"createdAt": u.CreatedAt,
	}
	writeJSON(w, resp, http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := contextWithTimeout(r, 5*time.Second)
	defer cancel()

	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := h.auth.Login(ctx, c.Email, c.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	writeJSON(w, map[string]string{"token": token}, http.StatusOK)
}

// Helper func
func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func contextWithTimeout(r *http.Request, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), d)
}
