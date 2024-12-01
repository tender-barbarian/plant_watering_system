package service

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/tender-barbarian/gniot/webserver/internal/repository"
)

// SensorService is the service to manage the Sensors.
type SensorService struct {
	SensorRepository       *repository.SensorRepository
	SensorMethodRepository *repository.SensorMethodRepository
}

// NewSensorService returns a new NewSensorService.
func NewSensorService(SensorRepository *repository.SensorRepository, SensorMethodRepository *repository.SensorMethodRepository) *SensorService {
	return &SensorService{
		SensorRepository:       SensorRepository,
		SensorMethodRepository: SensorMethodRepository,
	}
}

// List returns a list of all Sensors.
func (s *SensorService) List(ctx context.Context) ([]repository.SensorModel, error) {
	return s.SensorRepository.FindAll(ctx)
}

// List returns a list of all SensorMethods.
func (s *SensorService) ListMethods(ctx context.Context, sensorMethodIDs []int32) ([]repository.SensorMethodModel, error) {
	return s.SensorMethodRepository.FindAll(ctx, sensorMethodIDs)
}

// Create creates a new Sensor.
func (s *SensorService) Create(ctx context.Context, name string, sensorType string, chip string, board string) (int, error) {
	return s.SensorRepository.Create(ctx, repository.SensorRepositoryCreateParams{
		Name:       name,
		SensorType: sensorType,
		Chip:       chip,
		Board:      board,
	})
}

// Create creates a new SensorMethod.
func (s *SensorService) CreateMethod(ctx context.Context, name string, HttpMethod string, RequestBody string, board string) (int, error) {
	return s.SensorMethodRepository.Create(ctx, repository.SensorMethodRepositoryCreateParams{
		Name:        name,
		HttpMethod:  HttpMethod,
		RequestBody: RequestBody,
	})
}

// Get returns a Sensor by id.
func (s *SensorService) Get(ctx context.Context, id int) (repository.SensorModel, error) {
	return s.SensorRepository.Find(ctx, id)
}

// Get returns a SensorMethod by id.
func (s *SensorService) GetMethod(ctx context.Context, id int) (repository.SensorMethodModel, error) {
	return s.SensorMethodRepository.Find(ctx, id)
}

// Delete deletes a Sensor by id.
func (s *SensorService) Delete(ctx context.Context, id int) error {
	return s.SensorRepository.Delete(ctx, id)
}

// Delete deletes a SensorMethod by id.
func (s *SensorService) DeleteMethod(ctx context.Context, id int) error {
	return s.SensorMethodRepository.Delete(ctx, id)
}

func (s *SensorService) ExecuteMethod(ctx context.Context, sensorId int, methodName string) error {
	sensor, err := s.Get(ctx, sensorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("cannot find Sensor with id %d: %v", sensorId, err)
		}
		return err
	}

	sensorMethods, err := s.ListMethods(ctx, sensor.SensorMethodIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Sensor with id %d has no methods: %v", sensorId, err)
		}
		return err
	}

	for _, sensorMethod := range sensorMethods {
		if methodName == sensorMethod.Name {
			err := s.sendHTTP(ctx, sensorMethod.HttpMethod, sensor.IP, sensorMethod.RequestBody)
			if err != nil {
				return fmt.Errorf("sending http request: %v", err)
			}
		}
	}

	return nil
}

func (s *SensorService) sendHTTP(ctx context.Context, method string, ip string, body string) error {
	r, err := http.NewRequestWithContext(ctx, method, ip, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
