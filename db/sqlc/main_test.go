package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
var testDB *sql.DB

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

func getConnectionString(containerType string, ctx context.Context) (string, error) {
	switch containerType {
	case "local":
		return dbSource, nil
	case "test":
		return createTestContainer(ctx)
	default:
		return "", errors.New("Coudn't get the connection right now!")
	}
}

func TestMain(m *testing.M) {
	ctx = context.Background()

	containerType := "test"

	connString, err := getConnectionString(containerType, ctx)
	fmt.Println("connection string", connString)
	if err != nil {
		log.Fatal("cannot create test container:", err)
	}

	testDB, err = sql.Open("postgres", connString)
	testDB.SetMaxOpenConns(10)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	code := m.Run()

	if containerType == "test" {
		container.Terminate(ctx)
	}

	os.Exit(code)
}
