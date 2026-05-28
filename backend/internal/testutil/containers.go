package testutil

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupPostgres(ctx context.Context) (func(), error) {
	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get postgres host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get postgres port: %w", err)
	}

	os.Setenv("POSTGRES_HOST", host)
	os.Setenv("POSTGRES_PORT", port.Port())
	os.Setenv("POSTGRES_DB", "test")
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASSWORD", "test")

	cleanup := func() {
		container.Terminate(ctx)
	}

	return cleanup, nil
}

func SetupMysql(ctx context.Context) (func(), error) {
	container, err := mysql.Run(ctx,
		"mysql:8.2",
		mysql.WithDatabase("test"),
		mysql.WithUsername("test"),
		mysql.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("3306/tcp").WithStartupTimeout(120*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start mysql container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mysql host: %w", err)
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mysql port: %w", err)
	}

	os.Setenv("MYSQL_HOST", host)
	os.Setenv("MYSQL_PORT", port.Port())
	os.Setenv("MYSQL_DATABASE", "test")
	os.Setenv("MYSQL_USER", "test")
	os.Setenv("MYSQL_PASSWORD", "test")

	time.Sleep(5 * time.Second)

	cleanup := func() {
		container.Terminate(ctx)
	}

	return cleanup, nil
}
