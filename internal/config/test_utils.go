package config

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupTestDB initializes the test database connection and returns a cleanup function.
func SetupTestDB(m *testing.M) func() {
	connectionString := "postgres://postgres:password@127.0.0.1:5432/jwt_auth_test?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Panic("failed to connect to test database: " + err.Error())
	}

	DbConn = &DbConnect{pool: pool}

	// Clean tables before running tests
	_, err = DbConn.GetPool().Exec(context.Background(), "TRUNCATE TABLE refresh_tokens, users RESTART IDENTITY CASCADE")
	if err != nil {
		panic("failed to truncate tables: " + err.Error())
	}

	return func() {
		DbConn.Close()
	}
}
