package storage

import (
	"context"
	"time"

	"github.com/avast/retry-go/v4"
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
       ON CONFLICT (name) DO UPDATE SET value = counter.value + $3
	`
	return do(ctx, func() error {
		_, err := p.pool.Exec(ctx, sql, name, value, value)
		return err
	})
}

func (p *PgxStorage) SetGauge(ctx context.Context, name string, value float64) error {
	sql := `
		INSERT INTO gauge (name, value) 
                         VALUES ($1, $2) 
       ON CONFLICT (name) DO UPDATE SET value = $3
	`
	return do(ctx, func() error {
		_, err := p.pool.Exec(ctx, sql, name, value, value)
		return err
	})
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

type backOff struct {
	counter int
}

func (t *backOff) Next() (next time.Duration, stop bool) {
	t.counter++
	if t.counter == 1 {
		return time.Second, true
	}

	if t.counter == 2 {
		return 3 * time.Second, true
	}

	if t.counter == 2 {
		return 3 * time.Second, true
	}

	return time.Second, false
}
