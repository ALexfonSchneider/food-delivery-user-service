package postgres

import (
	"database/sql"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/config"
	"github.com/pressly/goose/v3"
)

func (r *Repository) Migrate(cfg *config.Config) error {
	db, err := sql.Open("postgres", cfg.Postgres.ConnectionStringPQ())

	if err != nil {
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
