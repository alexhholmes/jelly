package healthcheck

import (
	"log/slog"
	"net/http"

	"jelly/pkg/api/v1/gen"
	util2 "jelly/pkg/api/v1/util"
)

// HealthHandler implements health check endpoints.
type HealthHandler struct{}

// HealthCheck returns {"status": "ok"} with HTTP 200.
// GET /health
func (h HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(util2.ContextLogger).(*slog.Logger)
	logger.Info("Healthcheck")

	resp := gen.HealthCheck{
		Status: "ok",
	}

	util2.WriteJSONResponse(w, logger, http.StatusOK, resp)
}
