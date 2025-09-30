package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func RunMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("init migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations", // путь до миграций
		"postgres",                   // имя БД
		driver,
	)
	if err != nil {
		return fmt.Errorf("init migrate instance: %w", err)
	}

	if err := m.Drop(); err != nil {
		return fmt.Errorf("drop failed: %w", err)
	}

	driver, _ = postgres.WithInstance(db, &postgres.Config{})
	m, _ = migrate.NewWithDatabaseInstance(
		"file://internal/migrations",
		"postgres",
		driver,
	)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("up failed: %w", err)
	}

	return nil
}
