package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	cleanup := SetupTestDB(m)
	defer cleanup()
	os.Exit(m.Run())
}