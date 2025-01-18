package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgOnce sync.Once
	pgConn *pgxpool.Pool
)

func GetDbConnectionUrl(username string, password string, host string, port string, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
}

func GetDbConnectionUrlFromEnv() string {
	// postgres://[user]:[password]@[host]:[port]/[dbname]
	connString := GetDbConnectionUrl(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	log.Println("Connecting to:", connString)

	return connString
}

func GetPostgresConnection(connString string) (*pgxpool.Pool, error) {
	var err error

	pgOnce.Do(func() {
		pgConn, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatal("Error to create postgres database connection:", err)
		}

		var testResult int
		err = pgConn.QueryRow(context.Background(), "SELECT 1").Scan(&testResult)
		if err != nil {
			log.Fatal("Postgres failed to connect:", err)
		}

		log.Println("Postgres database connection successfully obtained")
	})

	return pgConn, err
}

func Setup(ctx context.Context) *pgxpool.Pool {
	log.SetPrefix("DB: ")

	pgConn, err := GetPostgresConnection(GetDbConnectionUrlFromEnv())
	if err != nil {
		log.Fatal("Error getting database connection:", err)
	}

	log.SetPrefix("")

	return pgConn
}
