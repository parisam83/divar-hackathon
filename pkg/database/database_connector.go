package database

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/url"
	"time"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/internal/configs"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFileSystem embed.FS


func buildConnectionURL(config configs.DatabaseConfig) (*url.URL, error) {
    pgURL := &url.URL{
        Scheme: "postgres",
        User:   url.UserPassword(config.Username, config.Password),
        Host:   fmt.Sprintf("%s:%d", config.Host, config.Port),
        Path:   "/" + config.DBName,
    }

    query := pgURL.Query()
    query.Add("sslmode", config.SSLMode)
    pgURL.RawQuery = query.Encode()
	
    return pgURL, nil
}

func configurePoolConfig(config configs.DatabaseConfig, connStr string) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(connStr)
    if err != nil {
		log.Printf("failed to parse postgres config URL: %v", err)
		return nil, err
	}

    poolConfig.MaxConns = config.MaxConns
    poolConfig.MinConns = config.MinConns
    poolConfig.MaxConnLifetime = time.Duration(config.MaxConnLifetimeMinutes) * time.Minute
    poolConfig.MaxConnIdleTime = time.Duration(config.MaxConnIdleTimeMinutes) * time.Minute
    poolConfig.MaxConnLifetimeJitter = time.Duration(config.MaxConnLifetimeJitterMinutes) * time.Minute

    return poolConfig, nil
}

func runMigrations(connURL string) error {
    driver, err := iofs.New(migrationFileSystem, "migrations")
    if err != nil {
		log.Printf("failed to load migrations: %v", err)
        return err
    }
	
    m, err := migrate.NewWithSourceInstance("iofs", driver, connURL)
    if err != nil {
		log.Printf("failed to create migration instance: %v", err)
		return err
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("failed to apply migrations: %v", err)
		return err
    }

    return nil
}

func ConnectToDatabase(serverContext context.Context, databaseConfig configs.DatabaseConfig) (*pgxpool.Pool, error) {
	pgURL, err := buildConnectionURL(databaseConfig)
	if err != nil {
		log.Printf("failed to build postgres connection URL: %v", err)
		return nil, err
	}

	connStr := pgURL.String()
	pgxPoolConfig, err := configurePoolConfig(databaseConfig, connStr)
	if err != nil {
		log.Printf("failed to configure postgres pool config: %v", err)
		return nil, err
	}

	// Establish Connection Pool
	pool, err := pgxpool.NewWithConfig(serverContext, pgxPoolConfig)
	if err != nil {
		log.Printf("failed to create PGX connection pool: %v", err)
		return nil, err	
	}

	// Verify connection
	if err := pool.Ping(serverContext); err != nil {
		pool.Close()
		log.Printf("failed to ping database: %v", err)
		return nil, err
	}

	// Run migrations
	if err := runMigrations(pgURL.String()); err != nil {
		pool.Close()
		log.Printf("failed to run migrations: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database and applied migrations")
	return pool, nil
}