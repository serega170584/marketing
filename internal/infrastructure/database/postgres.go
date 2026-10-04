package database

import (
	"context"
	"fmt"
	"marketing/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

//nolint:gosec
func New(cfg config.DBConfig) (*Postgres, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxOpenConns)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *Postgres) Pool() *pgxpool.Pool {
	return p.pool
}
