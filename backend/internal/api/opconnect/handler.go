package opconnect

import (
	"net/http"
	"time"

	"github.com/shahnajsc/OnePointLedger/backend/internal/api/middleware"
	"github.com/shahnajsc/OnePointLedger/backend/internal/service"
)

type Handler struct {
	svc *service.OPConnectService
}

func NewHandler(svc *service.OPConnectService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		http.Error(w, "missing user context", http.StatusUnauthorized)
		return
	}

	ctx, cancel := contextWithTimeout(r, 15*time.Second)
	defer cancel()

	authURL, err := h.svc.Start(ctx, userID)
	if err != nil {
		http.Error(w, "failed to start OP connect: "+err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]string{"authorization_url": authURL}, http.StatusOK)
}
