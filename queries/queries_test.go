package queries

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewQueries(db *sql.DB) *Queries {
	return &Queries{
		Database: db,
	}
}

func TestGetLogsByElevator_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := "SELECT id, elevator_id, current_floor, state, direction, timestamp, query FROM logs WHERE elevator_id = $1"
	mock.ExpectQuery(query).WithArgs(1).WillReturnError(sql.ErrConnDone)

	queriesInstance := NewQueries(db)

	logs, err := queriesInstance.GetLogsByElevator(1)
	assert.Error(t, err)
	assert.Nil(t, logs)
}

func TestGetBuilding(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := "SELECT id, no_of_floors, name FROM buildings WHERE id = $1"
	mockRows := sqlmock.NewRows([]string{"id", "no_of_floors", "name"}).
		AddRow(1, 5, "Building A")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(mockRows)

	queriesInstance := NewQueries(db)

	building, err := queriesInstance.GetBuilding(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, building.ID)
	assert.Equal(t, 5, building.TotalFloors)
	assert.Equal(t, "Building A", building.Name)
}

func TestUpdateElevatorState(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := "UPDATE elevators SET floor = $1, direction = $2, moving = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4"
	mock.ExpectExec(query).WithArgs(2, "up", true, 1).WillReturnResult(sqlmock.NewResult(0, 1))

	queriesInstance := NewQueries(db)

	err := queriesInstance.UpdateElevatorState(1, 2, "up", true)
	assert.NoError(t, err)
}

func TestInsertLogs(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	insertQuery := `
		INSERT INTO logs (elevator_id, current_floor, state, direction, timestamp, query)
		VALUES
			($1, $2, $3, $4, CURRENT_TIMESTAMP, $5);
	`
	mock.ExpectExec(insertQuery).WithArgs(1, 5, "door-open", "up", "encoded_query").WillReturnResult(sqlmock.NewResult(0, 1))

	queriesInstance := NewQueries(db)

	err := queriesInstance.InsertLogs(1, 5, "door-open", "up", true)
	assert.NoError(t, err)
}
