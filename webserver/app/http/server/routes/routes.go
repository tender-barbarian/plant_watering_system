package routes

import (
	"context"
	"net/http"

	"github.com/tender-barbarian/gniot/webserver/app/http/server/handlers"
)

type Routes struct {
	handlers *handlers.Handlers
}

func NewRoutes(handlers *handlers.Handlers) *Routes {
	return &Routes{handlers: handlers}
}

func (r *Routes) Add(ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/sensor/{id}", r.handlers.GetSensor(ctx))
	mux.HandleFunc("/sensor/{id}/{sensor_method}", r.handlers.ExecSensorMethod(ctx))
	mux.Handle("/", http.NotFoundHandler())

	return mux
}
