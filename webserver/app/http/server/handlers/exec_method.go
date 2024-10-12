package handlers

import (
	"context"
	"net/http"
	"strconv"
)

// Handle handles the http request.
func (h *Handlers) ExecSensorMethod(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sensorId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		methodName := r.PathValue("sensor_method")

		err = h.service.ExecuteMethod(ctx, sensorId, methodName)
		if err != nil {
			h.errorsWrapper.ServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
