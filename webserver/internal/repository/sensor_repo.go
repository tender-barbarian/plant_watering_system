package repository

import (
	"context"
	"database/sql"
	"sync"

	sq "github.com/Masterminds/squirrel"
)

type SensorModel struct {
	ID              int32   `json:"id"`
	Name            string  `json:"name"`
	SensorType      string  `json:"sensorType"`
	Chip            string  `json:"chip"`
	Board           string  `json:"board"`
	IP              string  `json:"ip"`
	SensorMethodIDs []int32 `json:"sensorMethodIDs,omitempty"`
}

// SensorRepository is the repository to handle the [repository.Sensor] repository database interactions.
type SensorRepository struct {
	mutex sync.Mutex
	db    *sql.DB
}

// NewSensorRepository returns a new [SensorRepository].
func NewSensorRepository(db *sql.DB) *SensorRepository {
	return &SensorRepository{
		db: db,
	}
}

// Find finds a Sensor by id.
func (r *SensorRepository) Find(ctx context.Context, id int) (SensorModel, error) {
	var sensor SensorModel

	query, args, err := sq.
		Select("id", "name", "sensorType", "Chip", "Board", "SensorMethodIDs").
		From("Sensors").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return sensor, err
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	err = row.Scan(&sensor.ID, &sensor.Name, &sensor.SensorType, &sensor.Chip, &sensor.Board, &sensor.SensorMethodIDs)

	return sensor, err
}

// FindAll finds all Sensors.
func (r *SensorRepository) FindAll(ctx context.Context) ([]SensorModel, error) {
	qb := sq.
		Select("id", "name", "sensorType", "chip", "board", "SensorMethodIDs").
		From("Sensors").
		OrderBy("id")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var sensors []SensorModel
	for rows.Next() {
		var sensor SensorModel
		if err = rows.Scan(&sensor.ID, &sensor.Name, &sensor.SensorType, &sensor.Chip, &sensor.Board, &sensor.SensorMethodIDs); err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}

// SensorRepositoryCreateParams is a parameter for Create.
type SensorRepositoryCreateParams struct {
	Name            string
	SensorType      string
	Chip            string
	Board           string
	SensorMethodIDs []int32
}

// Create creates a new Sensor and returns its id.
func (r *SensorRepository) Create(ctx context.Context, params SensorRepositoryCreateParams) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query, args, err := sq.Insert("Sensors").Columns("name", "sensorType", "chip", "board", "sensorMethodIds").Values(params.Name, params.SensorType, params.Chip, params.Board, params.SensorMethodIDs).ToSql()
	if err != nil {
		return 0, err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Delete deletes an existing Sensor by id.
func (r *SensorRepository) Delete(ctx context.Context, id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query, args, err := sq.Delete("Sensors").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err
}
