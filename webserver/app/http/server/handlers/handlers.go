package handlers

import (
	"log/slog"

	"github.com/tender-barbarian/gniot/webserver/app/http/server/errors"
	"github.com/tender-barbarian/gniot/webserver/internal/service"
)

type Handlers struct {
	logger        *slog.Logger
	errorsWrapper errors.ErrorsWrapper
	service       *service.SensorService
}

func NewHandlers(sensorService *service.SensorService, errorsWrapper errors.ErrorsWrapper, logger *slog.Logger) *Handlers {
	return &Handlers{
		logger:        logger,
		errorsWrapper: errorsWrapper,
		service:       sensorService,
	}
}
