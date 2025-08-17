// CLI program that runs DB migrations (up or down) against a SQLite database using golang-migrate.
package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]

	// Open a database handle for the SQLite file in the working directory.
    // sql.Open does not create a file automatically for all drivers; this is the
    // standard way to obtain *sql.DB for migrate to use.
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Wrap the *sql.DB in the migrate sqlite3 driver instance.
	instance, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Create a file source pointing to the migrations directory.
    // Path is relative to the process working directory (run from repo root).
	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}

	// Create the migrate object using the file source and sqlite3 database instance.
	m, err := migrate.NewWithInstance("file", fSrc, "sqlite", instance)
	if err != nil {
		log.Fatal(err)
	}

	switch direction{
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange{
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange{
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid migration direction. Please use 'up' or 'down'.")
	}
}