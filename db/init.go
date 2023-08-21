package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func SeedTable(db *sql.DB) {
	buildingSeedSql := `
		INSERT INTO buildings (no_of_floors, name)
		VALUES
			(10, 'Office Building');
	`

	// Seed data for elevators table
	elevatorSeedSql := `
		INSERT INTO elevators (building_id, floor, direction, moving, is_door_open)
		VALUES
			(1, 1, 'none', false, false);
	`

	seedQueries := []string{
		buildingSeedSql,
		elevatorSeedSql,
	}

	for _, query := range seedQueries {
		_, err := db.Exec(query)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("✅ Seed function ran")
	}
}

func createTables(db *sql.DB) {
	createBuildingTableSql := `
		CREATE TABLE IF NOT EXISTS buildings (
			id SERIAL PRIMARY KEY,
			no_of_floors INT NOT NULL,
			name VARCHAR(255),
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`

	createElevatorTableSql := `
		CREATE TABLE IF NOT EXISTS elevators (
			id SERIAL PRIMARY KEY,
			building_id INT REFERENCES buildings (id),
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
			current_floor INT NOT NULL,
			state VARCHAR(10) NOT NULL CHECK (state IN ('door-open', 'door-closed', 'moving', 'stopped')),
			direction VARCHAR(10) NOT NULL CHECK (direction IN ('none', 'up', 'down')),
			timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP 
		);
	`

	alterLogTableSql := `
		ALTER TABLE logs
		ADD COLUMN IF NOT EXISTS query TEXT NOT NULL;
	`

	tableQueries := []string{
		createBuildingTableSql,
		createElevatorTableSql,
		createLogTableSql,
		alterLogTableSql,
	}

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

	connectionString := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ :::Connected to the database:::")
	createTables(db)

	return db
}
