package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDb *pgxpool.Pool

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error

	testDb, err = pgxpool.New(context.Background(), dbSource)

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")
	testQueries = New(testDb)
	os.Exit(m.Run())
}
