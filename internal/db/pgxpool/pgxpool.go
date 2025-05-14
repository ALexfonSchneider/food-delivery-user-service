package pgxpool

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/config"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func MustPGXPool(ctx context.Context, cfg *config.Config, logger *slog.Logger, level slog.Level) *pgxpool.Pool {
	conf, err := pgxpool.ParseConfig(cfg.Postgres.ConnectionString())
	if err != nil {
		panic(err)
	}

	conf.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	return pool
}
