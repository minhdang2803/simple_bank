package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/minhdang2803/simple_bank/utils"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
// )

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load configuration")
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database")
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
