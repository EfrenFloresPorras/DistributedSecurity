package http

import (
	"encoding/json"
	"net/http"

	"DistributedSecurity/auth-service/internal/controller/auth"
	"DistributedSecurity/auth-service/pkg/model"
)

type Handler struct {
	ctrl *auth.Controller
}

func New(ctrl *auth.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	resp, err := h.ctrl.Login(req, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
