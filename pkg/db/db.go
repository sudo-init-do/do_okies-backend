package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustOpen(ctx context.Context, url string) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil { panic(err) }
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.HealthCheckPeriod = 30 * time.Second
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil { panic(err) }
	if err := pool.Ping(ctx); err != nil { panic(err) }
	return pool
}
