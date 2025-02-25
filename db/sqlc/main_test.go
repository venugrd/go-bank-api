package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var container *postgres.PostgresContainer
var ctx context.Context

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func createTestContainer(ctx context.Context) (connString string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testScript := wd + "/testdata/schema.sql"
	container, err = postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithInitScripts(testScript),
		postgres.WithDatabase("simple_bank"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		panic(err)
	}

	connString, err = container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	return connString, nil
}

func TestMain(m *testing.M) {
	ctx = context.Background()

	connString, err := createTestContainer(ctx)
	if err != nil {
		log.Fatal("cannot create test container:", err)
	}

	conn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)

	code := m.Run()

	container.Terminate(ctx)

	os.Exit(code)
}
