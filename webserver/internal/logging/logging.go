package logging

import (
	"log/slog"
	"net/http"
)

type Logger struct {
	handler http.Handler
	logger  *slog.Logger
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ip     = r.RemoteAddr
		proto  = r.Proto
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	l.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
}

func NewLoggingMiddleware(handler http.Handler, logger *slog.Logger) *Logger {
	return &Logger{handler, logger}
}
