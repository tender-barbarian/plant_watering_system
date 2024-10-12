package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Handle handles the http request.
func (h *Handlers) GetSensor(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sensorId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		sensor, err := h.service.Get(ctx, sensorId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			h.errorsWrapper.ServerError(w, r, err)
			return
		}

		response, err := json.Marshal(sensor)
		if err != nil {
			h.errorsWrapper.ServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(response)
		if err != nil {
			h.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		}
	}
}
