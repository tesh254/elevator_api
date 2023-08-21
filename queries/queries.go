package queries

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Queries struct {
	Database *sql.DB
}

type Building struct {
	ID          int    `json:"id"`
	TotalFloors int    `json:"no_of_floors"`
	Name        string `json:"name,omitempty"`
}

type Elevator struct {
	ID          int     `json:"id"`
	BuildingID  int     `json:"building_id"`
	Floor       int     `json:"floor"`
	Direction   string  `json:"direction"`
	Moving      bool    `json:"moving"`
	IsDoorsOpen bool    `json:"is_door_open"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

type Log struct {
	ID           int       `json:"id"`
	ElevatorID   int       `json:"elevator_id"`
	CurrentFloor int       `json:"current_floor"`
	State        string    `json:"state"`
	Direction    string    `json:"direction"`
	Timestamp    time.Time `json:"timestamp"`
	Query        string    `json:"query"`
}

func currentTimeStamp() string {
	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	return formattedTime
}

func Encode(text string) string {
	data := []byte(text)

	encoded := base64.StdEncoding.EncodeToString(data)

	return encoded
}

func Decode(encodedText string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedText)

	if err != nil {
		return "", nil
	}

	return string(decoded), nil
}

func (p *Queries) GetElevatorById(id int) (*Elevator, error) {
	query := "SELECT * FROM elevators WHERE id = $1"
	row := p.Database.QueryRow(query, id)

	var elevator Elevator
	err := row.Scan(
		&elevator.ID,
		&elevator.BuildingID,
		&elevator.Floor,
		&elevator.Direction,
		&elevator.Moving,
		&elevator.IsDoorsOpen,
		&elevator.CreatedAt,
		&elevator.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &elevator, nil
}

func (p *Queries) GetBuilding(id int) (*Building, error) {
	query := "select id, no_of_floors, name from buildings where id = $1"
	row := p.Database.QueryRow(query, id)

	var building Building

	err := row.Scan(
		&building.ID, &building.TotalFloors, &building.Name,
	)

	if err != nil {
		return nil, err
	}

	return &building, nil
}

func (p *Queries) UpdateElevatorState(id int, floor int, direction string, moving bool) error {
	query := "update elevators set floor = $1, direction = $2, moving = $3, updated_at = CURRENT_TIMESTAMP where id = $4"
	_, err := p.Database.Exec(query, floor, string(direction), moving, id)

	if err != nil {
		return err
	}

	return nil
}

func (p *Queries) InsertLogs(elevatorId int, currentFloor int, state string, direction string, moving bool) error {
	defaultElevatorUpdateQuery := fmt.Sprintf("update elevators set floor = %d, direction = %s, moving = %t, updated_at = %s where id = %d", currentFloor, direction, moving, currentTimeStamp(), elevatorId)

	encryptedDefaultElevatorQuery := Encode(defaultElevatorUpdateQuery)

	insertQuery := `
		insert into logs (elevator_id, current_floor, state, direction, timestamp, query)
		values
			($1, $2, $3, $4, CURRENT_TIMESTAMP, $5);
	`

	_, err := p.Database.Exec(insertQuery, elevatorId, currentFloor, state, direction, string(encryptedDefaultElevatorQuery))

	if err != nil {
		return err
	}

	return nil
}

func (p *Queries) GetLogsByElevator(elevatorId int) (*[]Log, error) {
	var logs []Log

	query := "SELECT id, elevator_id, current_floor, state, direction, timestamp, query FROM logs WHERE elevator_id = $1"
	rows, err := p.Database.Query(query, elevatorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log Log
		if err := rows.Scan(&log.ID, &log.ElevatorID, &log.CurrentFloor, &log.State, &log.Direction, &log.Timestamp, &log.Query); err != nil {
			return nil, err
		}

		prevEncodedQuery := log.Query

		decodedQuery, _ := Decode(prevEncodedQuery)

		log.Query = decodedQuery

		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &logs, nil
}
