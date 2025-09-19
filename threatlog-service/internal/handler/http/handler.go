package http

import (
	"encoding/json"
	"net/http"

	"DistributedSecurity/threatlog-service/internal/controller/threat"
)

type Handler struct {
	ctrl *threat.Controller
}

func New(ctrl *threat.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// /healthz
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// /log
func (h *Handler) Log(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr // simplificado: IP del cliente
	event := h.ctrl.LogEvent("failed_login", ip)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// /events
func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	events := h.ctrl.ListEvents()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
