package user

import (
	"net/http"

	"github.com/shahnajsc/OnePointLedger/backend/internal/api/middleware"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You are user: " + userID))
}
