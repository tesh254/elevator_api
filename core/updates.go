package core

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"elevator/queries"
)

type Elevator struct {
	ID           int    `json:"id"`
	BuildingID   int    `json:"building_id"`
	CurrentFloor int    `json:"current_floor"`
	Direction    string `json:"direction"`
	Moving       bool   `json:"moving"`
	DoorsOpen    bool   `json:"doors_open"`
	ToFloor      int    `json:"to_floor"`
	WsConn       *websocket.Conn
	DbConn       *sql.DB
	Clients      map[*websocket.Conn]bool
}

func generateLog(eventType, description string) string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] Event: %s - Description: %s", currentTime, eventType, description)
}

func (e *Elevator) SaveLog() {
	var state string
	if e.Moving && e.CurrentFloor != e.ToFloor {
		state = "moving"
	}

	if e.Moving && e.CurrentFloor == e.ToFloor {
		state = "stopped"
	}

	if e.DoorsOpen && !e.Moving {
		state = "door-open"
	}

	if !e.DoorsOpen && !e.Moving {
		state = "door-closed"
	}

	var queries queries.Queries

	queries.Database = e.DbConn

	err := queries.InsertLogs(e.ID, e.CurrentFloor, state, e.Direction, e.Moving)

	if err != nil {
		fmt.Println(err)
	}
}

func (e *Elevator) SendUpdatesToClients(message string) {
	e.SaveLog()
	for client := range e.Clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("WebSocket write error:", err)
			delete(e.Clients, client)
			break
		}
	}
}

func (e *Elevator) MoveToFloor(floor int) {
	var queries queries.Queries

	queries.Database = e.DbConn

	if e.ToFloor > e.CurrentFloor {
		e.Direction = "up"
	} else if floor < e.CurrentFloor {
		e.Direction = "down"
	} else {
		e.Direction = "none"
	}
	queries.UpdateElevatorState(e.ID, e.CurrentFloor, e.Direction, e.Moving)
	e.SendUpdatesToClients(generateLog("[update]", fmt.Sprintf("Elevator %d direction is set to %s", e.ID, e.Direction)))

	e.Moving = true
	queries.UpdateElevatorState(e.ID, e.CurrentFloor, e.Direction, e.Moving)

	for e.CurrentFloor != e.ToFloor {
		if e.Direction == "up" {
			e.CurrentFloor++
		} else if e.Direction == "down" {
			e.CurrentFloor--
		}
		queries.UpdateElevatorState(e.ID, e.CurrentFloor, e.Direction, e.Moving)
		e.SendUpdatesToClients(generateLog("[update]", fmt.Sprintf("Elevator %d direction is set to %s, moving state: %t", e.ID, e.Direction, e.Moving)))
		time.Sleep(5 * time.Second)
	}
	e.Moving = false
	queries.UpdateElevatorState(e.ID, e.CurrentFloor, e.Direction, e.Moving)
}

func (e *Elevator) OpenDoorsOnCorrectFloor(floor int) {
	var queries queries.Queries
	queries.Database = e.DbConn

	if e.CurrentFloor == e.ToFloor {
		e.DoorsOpen = true
		e.SendUpdatesToClients(generateLog("[update]", fmt.Sprintf("Elevator %d doors are open", e.ID)))
		time.Sleep(2 * time.Second)
		e.DoorsOpen = false
		queries.UpdateElevatorState(e.ID, e.CurrentFloor, "none", e.Moving)
		e.SendUpdatesToClients(generateLog("[update]", fmt.Sprintf("Elevator %d doors are closed", e.ID)))
	}
}

func (e *Elevator) Start() {
	var queries queries.Queries

	queries.Database = e.DbConn

	building, err := queries.GetBuilding(e.BuildingID)

	if err != nil {
		fmt.Println(err)
		e.SendUpdatesToClients(generateLog("[error]", "Failed to get building"))
	}

	numberOfFloors := building.TotalFloors
	var elevator Elevator = *e

	destinationFloor := e.ToFloor

	go func(e Elevator) {
		for {
			if e.Moving {
				continue
			}

			if destinationFloor == e.CurrentFloor {
				break
			}

			e.MoveToFloor(destinationFloor)
			e.SendUpdatesToClients(generateLog("[update]", fmt.Sprintf("Elevator %d reached destination floor %d of %d building floors\n", e.ID, e.CurrentFloor, numberOfFloors)))
			e.OpenDoorsOnCorrectFloor(destinationFloor)
		}
	}(elevator)
}
