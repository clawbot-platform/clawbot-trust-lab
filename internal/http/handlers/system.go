package handlers

import (
	"context"
	"net/http"

	"clawbot-trust-lab/internal/version"
)

type SystemHandler struct {
	readiness func(context.Context) error
	info      version.Info
}

func NewSystemHandler(readiness func(context.Context) error, info version.Info) *SystemHandler {
	return &SystemHandler{readiness: readiness, info: info}
}

func (h *SystemHandler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *SystemHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if err := h.readiness(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "not_ready",
			"error":  err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ready"})
}

func (h *SystemHandler) Version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.info)
}
