package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type Database struct {
	DSN       string
	container *postgres.PostgresContainer
}

func Start(ctx context.Context) (*Database, error) {
	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("gestao_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("start PostgreSQL test container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, terminateAfterError(container, fmt.Errorf("get PostgreSQL test connection string: %w", err))
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, terminateAfterError(container, fmt.Errorf("open PostgreSQL test connection: %w", err))
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, terminateAfterError(container, closeAfterError(db, fmt.Errorf("ping PostgreSQL test container: %w", err)))
	}
	if err := applyMigrations(ctx, db); err != nil {
		return nil, terminateAfterError(container, closeAfterError(db, err))
	}
	if err := db.Close(); err != nil {
		return nil, terminateAfterError(container, fmt.Errorf("close PostgreSQL migration connection: %w", err))
	}

	return &Database{DSN: dsn, container: container}, nil
}

func (db *Database) Terminate(ctx context.Context) error {
	if err := db.container.Terminate(ctx); err != nil {
		return fmt.Errorf("terminate PostgreSQL test container: %w", err)
	}
	return nil
}

func applyMigrations(ctx context.Context, db *sql.DB) error {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("locate test database migrations")
	}
	migrationsDir := filepath.Join(filepath.Dir(sourceFile), "..", "..", "shared", "plataforma", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read PostgreSQL migrations: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		migration, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}
		if _, err := db.ExecContext(ctx, string(migration)); err != nil {
			return fmt.Errorf("apply migration %s: %w", entry.Name(), err)
		}
	}
	return nil
}

func closeAfterError(db *sql.DB, cause error) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("%w; also failed to close PostgreSQL migration connection: %v", cause, err)
	}
	return cause
}

func terminateAfterError(container *postgres.PostgresContainer, cause error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := container.Terminate(ctx); err != nil {
		return fmt.Errorf("%w; also failed to terminate test container: %v", cause, err)
	}
	return cause
}
