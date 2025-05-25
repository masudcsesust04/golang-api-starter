package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps the pgxpool.Pool connection pool
type DbConnect struct {
	pool *pgxpool.Pool
}

var DbConn *DbConnect

// NewDB creates a new database connection using pgxpool
func InitDB(databaseURL string) error {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Ping the database to verify the connection.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DbConn = &DbConnect{pool: pool}
	fmt.Println("Database connection established successfully!")

	return nil
}

// Close closes the database connection
func (db *DbConnect) Close() {
	db.pool.Close()
}

// GetPool returns the underlying pgxpool.Pool
func (db *DbConnect) GetPool() *pgxpool.Pool {
	// Check if the pool is nil or closed
	if db == nil || db.pool == nil {
		log.Panic("Database pool is not initialized")
		return nil
	}
	if err := db.pool.Ping(context.Background()); err != nil {
		fmt.Println("Database pool is not available:", err)
	}

	// Return the pool if it's available
	return db.pool
}
