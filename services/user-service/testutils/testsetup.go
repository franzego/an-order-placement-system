package testutils

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	godotenv.Load()
	dbUrl := os.Getenv("dbSource")
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}
	t.Cleanup(func() {
		conn.Close()
	})
	return conn
}
func CleanupTestData(t *testing.T, conn *pgxpool.Pool) {
	_, err := conn.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		t.Fatal("Failed to cleanup test data:", err)
	}
}
