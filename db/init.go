package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func seedTable(db *sql.DB) {
	elevatorSeedSql := `
		-- Seed data for elevators table
		INSERT INTO elevators (floor, direction, moving, is_door_open)
		VALUES
			(1, 'none', false, false);
	`

	var seedQueries [1]string

	seedQueries[0] = elevatorSeedSql

	for _, query := range seedQueries {
		_, err := db.Exec(query)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("✅ :::Seed function ran:::")
	}
}

func createTables(db *sql.DB) {
	createElevatorTableSql := `
		CREATE TABLE IF NOT EXISTS elevators (
			id SERIAL PRIMARY KEY,
			floor INT NOT NULL,
			direction VARCHAR(10) NOT NULL CHECK (direction IN ('none', 'up', 'down')),
			moving BOOLEAN NOT NULL,
			is_door_open BOOLEAN NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    		updated_at TIMESTAMPTZ
		)
	`

	createLogTableSql := `
		CREATE TABLE IF NOT EXISTS logs (
			id SERIAL PRIMARY KEY,
			elevator_id INT REFERENCES elevators (id),
			timestamp TIMESTAMPTZ NOT NULL,
			method VARCHAR(10) NOT NULL,
			path VARCHAR(255) NOT NULL,
			query TEXT,
			body TEXT,
			ip_address VARCHAR(45)
		);
	`

	var tableQueries [2]string

	tableQueries[0] = createElevatorTableSql
	tableQueries[1] = createLogTableSql

	for _, query := range tableQueries {
		_, err := db.Exec(query)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("✅ :::Tables created successfully:::")
	}
}

func ConnectToDatabase() *sql.DB {
	var db *sql.DB
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ :::Connected to the database:::")
	createTables(db)
	seedTable(db)

	return db
}
