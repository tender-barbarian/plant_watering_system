package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tender-barbarian/gniot/webserver/app/http/server/errors"
	"github.com/tender-barbarian/gniot/webserver/app/http/server/handlers"
	"github.com/tender-barbarian/gniot/webserver/app/http/server/routes"
	"github.com/tender-barbarian/gniot/webserver/internal/logging"
	"github.com/tender-barbarian/gniot/webserver/internal/repository"
	"github.com/tender-barbarian/gniot/webserver/internal/service"
)

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// Initialize repositories
	sensorRepository := repository.NewSensorRepository(nil)
	sensorMethodRepository := repository.NewSensorMethodRepository(nil)
	// Initalize service
	sensorService := service.NewSensorService(sensorRepository, sensorMethodRepository)
	// Initialize helpers
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	errorsWrapper := errors.NewErrorsWrapper(logger)
	// Initialize handlers and routes
	handlers := handlers.NewHandlers(sensorService, errorsWrapper, logger)
	routes := routes.NewRoutes(handlers)
	mux := routes.Add(ctx)
	// Initialize middleware
	var wrappedMux http.Handler = mux
	wrappedMux = logging.NewLoggingMiddleware(wrappedMux, logger)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("127.0.0.1", "80"),
		Handler: wrappedMux,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "listening and serving requests: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}
