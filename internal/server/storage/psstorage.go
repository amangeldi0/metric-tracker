package storage

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const schema = `
	CREATE TABLE IF NOT EXISTS metrics (
	    "id" SERIAL,
		"name" TEXT NOT NULL,
		"mtype" VARCHAR(12) NOT NULL DEFAULT 'gauge',
		"delta" INTEGER NOT NULL DEFAULT 0,
		"value" DOUBLE PRECISION NOT NULL DEFAULT 0.0,
		CONSTRAINT unique_id_mtype UNIQUE (name, mtype),
		PRIMARY KEY (id)
	)
`

type PsStorage struct {
	pgx    *pgxpool.Pool
	logger logger
}

func NewPsStorage(pgx *pgxpool.Pool, logger logger) (Storage, error) {
	dbStorage := &PsStorage{
		pgx:    pgx,
		logger: logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := dbStorage.pgx.Ping(ctx); err != nil {
		return nil, err
	}

	if _, err := dbStorage.pgx.Exec(ctx, schema); err != nil {
		return nil, err
	}

	return dbStorage, nil
}

func (dbStorage *PsStorage) SetGauge(name string, value float64) error {
	_, err := dbStorage.pgx.Exec(context.Background(), `
		INSERT INTO metrics (name, mtype, delta, value)
		VALUES ($1, 'gauge', 0, $2)
		ON CONFLICT (name, mtype) DO UPDATE 
		SET value = EXCLUDED.value
	`, name, value)
	return err
}

func (dbStorage *PsStorage) AddCounter(name string, value int64) error {
	_, err := dbStorage.pgx.Exec(context.Background(), `
		INSERT INTO metrics (name, mtype, delta, value)
		VALUES ($1, 'counter', $2, 0)
		ON CONFLICT (name, mtype) DO UPDATE 
		SET delta = metrics.delta + EXCLUDED.delta
	`, name, value)
	return err
}

func (dbStorage *PsStorage) GetGauge(name string) (float64, error) {
	var value float64
	err := dbStorage.pgx.QueryRow(context.Background(), `
		SELECT value FROM metrics WHERE name = $1 AND mtype = 'gauge'
	`, name).Scan(&value)
	return value, err
}

func (dbStorage *PsStorage) GetCounter(name string) (int64, error) {
	var value int64
	err := dbStorage.pgx.QueryRow(context.Background(), `
		SELECT delta FROM metrics WHERE name = $1 AND mtype = 'counter'
	`, name).Scan(&value)
	return value, err
}

func (dbStorage *PsStorage) GetAll() ([]MetricsValue, error) {
	rows, err := dbStorage.pgx.Query(context.Background(), `SELECT name, mtype, delta, value FROM metrics`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []MetricsValue
	for rows.Next() {
		var m MetricsValue
		if err := rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (dbStorage *PsStorage) Ping(ctx context.Context) error {
	return dbStorage.pgx.Ping(ctx)
}

func (dbStorage *PsStorage) Close() error {
	dbStorage.pgx.Close()
	return nil
}

func (dbStorage *PsStorage) GetMiddleware() gin.HandlerFunc {
	return func(_ *gin.Context) {}
}
