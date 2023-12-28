package test

import (
	"database/sql"
	"fmt"
	db "github/tdadadavid/fingreat/db/sqlc"
	"github/tdadadavid/fingreat/utils"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var testQuery *db.Store

const testDbName = "fingreat_test_db"
const sslmode = "?sslmode=disable"

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")

	if err != nil {
		log.Fatal("Could not load config")
	}

	// connct to postgres server.
  postgresConn, err := sql.Open(config.DbDriver, config.DBSource + sslmode);
	if err != nil {
		log.Fatalf("Could not connect to %s server, Error: %v", config.DbDriver, err)
	}

	// create database for testing.
	_, err = postgresConn.Exec(fmt.Sprintf("CREATE DATABASE %s;", testDbName))
	if err != nil {
		tearDown(postgresConn)
		log.Fatalf("Error creating test database %v", err)
	}

	// connect to the test database just created
	testDbConn, err := sql.Open(config.DbDriver, config.DBSource+sslmode);
	if err != nil {
		tearDown(postgresConn)
		log.Fatalf("Error creating test database %v", err)
	}

	// create driver to execute migration.
	driver, err := postgres.WithInstance(testDbConn, &postgres.Config{}); 
	if err != nil {
		tearDown(postgresConn)
		log.Fatalf("Error: cannot create migration driver %v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%v/", "../migrations"),
		config.DbDriver,
		driver,
	);
	if err != nil {
		tearDown(postgresConn)
		log.Fatalf("Error: migration setup failed %v", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		tearDown(postgresConn)
		log.Fatalf("Error: migrating up failed %v", err)
	}

	// initiate the testdb connection for the whoele
	// tests suite.
	testQuery = db.NewStore(testDbConn)

	// run the tests
	code := m.Run()

	// close connection to the test testdb
	testDbConn.Close();

	// close connection to postgres server.
	tearDown(postgresConn)

	// exit the process.
	os.Exit(code)
}

func tearDown(conn *sql.DB) {
	_, err := conn.Exec(fmt.Sprintf("DROP DATABASE %v WITH (FORCE)", testDbName));
	if err != nil {
		log.Fatalf("Error dropping database: %v", err)
	}
	conn.Close();
}