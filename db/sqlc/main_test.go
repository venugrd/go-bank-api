package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var container testcontainers.Container
var ctx context.Context

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func createTestContainer(ctx context.Context) (connString string, err error) {
	wd, err := os.Getwd()
	sqlScripts := wd + "/testdata/schema.sql"

	if err != nil {
		log.Fatal("Couldn't get current working deirectory", err)
		return "", err
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "simple_bank",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err = testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)

	if err != nil {
		log.Fatal("cannot create test container:", err)
		return "", fmt.Errorf("run container: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")

	if err != nil {
		log.Fatal(err)
	}

	connString = fmt.Sprintf("postgres://postgres:postgres@127.0.0.1:%d/simple_bank?sslmode=disable", port.Int())

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

	os.Exit(m.Run())
}
