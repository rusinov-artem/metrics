package storage

import (
	"context"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxStorage struct {
	pool  *pgxpool.Pool
	batch *pgx.Batch
}

func NewPgxStorage(pool *pgxpool.Pool) *PgxStorage {
	pgxStorage := &PgxStorage{
		pool:  pool,
		batch: &pgx.Batch{},
	}
	return pgxStorage
}

func (p *PgxStorage) SetCounter(name string, value int64) error {
	sql := `
		INSERT INTO counter (name, value) 
                         VALUES ($1, $2) 
       ON CONFLICT (name) DO UPDATE SET value = counter.value + excluded.value
	`
	p.batch.Queue(sql, name, value)

	return nil
}

func (p *PgxStorage) SetGauge(name string, value float64) error {
	sql := `
		INSERT INTO gauge (name, value) 
                         VALUES ($1, $2) 
       ON CONFLICT (name) DO UPDATE SET value = excluded.value
	`
	p.batch.Queue(sql, name, value)

	return nil
}

func (p *PgxStorage) Flush(ctx context.Context) error {
	result := p.pool.SendBatch(ctx, p.batch)
	err := result.Close()
	p.batch = &pgx.Batch{}
	return err
}

func (p *PgxStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	sql := `
		SELECT value from counter WHERE name = $1
	`
	var value int64
	err := do(ctx, func() error {
		row := p.pool.QueryRow(ctx, sql, name)
		return row.Scan(&value)
	})

	return value, err
}

func (p *PgxStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	sql := `
		SELECT value from gauge WHERE name = $1
	`

	var value float64
	err := do(ctx, func() error {
		row := p.pool.QueryRow(ctx, sql, name)
		return row.Scan(&value)
	})

	return value, err
}

func do(ctx context.Context, fn func() error) error {
	return retry.Do(
		fn,
		retry.Context(ctx),
		retry.Attempts(3),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			if n == 1 {
				return time.Second
			}
			if n == 2 {
				return 3 * time.Second
			}
			if n == 3 {
				return 5 * time.Second
			}
			return 10 * time.Second
		}),
	)
}
