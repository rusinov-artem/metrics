package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxStorage struct {
	pool *pgxpool.Pool
}

func NewPgxStorage(pool *pgxpool.Pool) *PgxStorage {
	return &PgxStorage{
		pool: pool,
	}
}

func (p *PgxStorage) SetCounter(ctx context.Context, name string, value int64) error {
	sql := `
		INSERT INTO counter (name, value) 
                         VALUES ($1, $2) 
       ON CONFLICT (name) DO UPDATE SET value = $3
	`
	_, err := p.pool.Exec(ctx, sql, name, value, value)
	return err
}

func (p *PgxStorage) SetGauge(ctx context.Context, name string, value float64) error {
	sql := `
		INSERT INTO couge (name, value) 
                         VALUES ($1, $2) 
       ON CONFLICT (name) DO UPDATE SET value = $3
	`
	_, err := p.pool.Exec(ctx, sql, name, value, value)
	return err
}

func (p *PgxStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	sql := `
		SELECT value from counter WHERE name = $1
	`
	row := p.pool.QueryRow(ctx, sql, name)
	var value int64
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (p *PgxStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	sql := `
		SELECT value from gauge WHERE name = $1
	`
	row := p.pool.QueryRow(ctx, sql, name)
	var value float64
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}

	return value, nil
}
