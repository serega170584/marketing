package database

import (
	"context"
	"fmt"
	"marketing/internal/config"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

//nolint:gosec
func New(
	ctx context.Context,
	cfg config.DBConfig,
) (*Postgres, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pgxTracer := otelpgx.NewTracer(
		otelpgx.WithIncludeQueryParameters(),
	)

	poolCfg.ConnConfig.Tracer = pgxTracer

	poolCfg.MaxConns = int32(cfg.MaxOpenConns)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
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
