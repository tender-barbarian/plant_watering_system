package repository

import (
	"context"
	"database/sql"
	"sync"

	sq "github.com/Masterminds/squirrel"
)

type SensorMethodModel struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	HttpMethod  string `json:"httpMehtod"`
	RequestBody string `json:"requestBody"`
}

// SensorMethodRepository is the Repository to handle the [SensorMethodModel] Repository database interactions.
type SensorMethodRepository struct {
	mutex sync.Mutex
	db    *sql.DB
}

// NewSensorMethodRepository returns a new [SensorMethodRepository].
func NewSensorMethodRepository(db *sql.DB) *SensorMethodRepository {
	return &SensorMethodRepository{
		db: db,
	}
}

// Find finds a SensorMethod by id.
func (r *SensorMethodRepository) Find(ctx context.Context, id int) (SensorMethodModel, error) {
	var sensorMethod SensorMethodModel

	query, args, err := sq.
		Select("id", "name", "httpMethod", "requestBody").
		From("SensorMethods").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return sensorMethod, err
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	err = row.Scan(&sensorMethod.ID, &sensorMethod.Name, &sensorMethod.HttpMethod, &sensorMethod.RequestBody)

	return sensorMethod, err
}

// FindAll finds all SensorMethods.
func (r *SensorMethodRepository) FindAll(ctx context.Context, sensorMethodIDs []int32) ([]SensorMethodModel, error) {
	qb := sq.
		Select("id", "name", "httpMethod", "requestBody").
		From("SensorMethods").
		Where("id IN ?", sensorMethodIDs).
		OrderBy("id")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var sensorMethods []SensorMethodModel
	for rows.Next() {
		var sensorMethod SensorMethodModel
		if err = rows.Scan(&sensorMethod.ID, &sensorMethod.Name, &sensorMethod.HttpMethod, &sensorMethod.RequestBody); err != nil {
			return nil, err
		}
		sensorMethods = append(sensorMethods, sensorMethod)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sensorMethods, nil
}

// SensorMethodRepositoryCreateParams is a parameter for Create.
type SensorMethodRepositoryCreateParams struct {
	Name        string
	HttpMethod  string
	RequestBody string
}

// Create creates a new SensorMethod and returns its id.
func (r *SensorMethodRepository) Create(ctx context.Context, params SensorMethodRepositoryCreateParams) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query, args, err := sq.Insert("SensorMethods").Columns("name", "httpMethod", "requestBody").Values(params.Name, params.HttpMethod, params.RequestBody).ToSql()
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

// Delete deletes an existing SensorMethod by id.
func (r *SensorMethodRepository) Delete(ctx context.Context, id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	query, args, err := sq.Delete("SensorMethods").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err
}
