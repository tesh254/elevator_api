package core

import (
	"time"
)

type ElevatorDirection string

const (
	None ElevatorDirection = "none"
	Up   ElevatorDirection = "up"
	Down ElevatorDirection = "down"
)

type Elevator struct {
	ID           int
	CurrentFloor int
	Direction    ElevatorDirection
	Moving       bool
	DoorsOpen    bool
}

func (e *Elevator) MoveToFloor(floor int) {
	if floor > e.CurrentFloor {
		e.Direction = "up"
	} else if floor < e.CurrentFloor {
		e.Direction = "down"
	} else {
		e.Direction = "none"
	}

	e.Moving = true
	for e.CurrentFloor != floor {
		if e.Direction == "up" {
			e.CurrentFloor++
		} else if e.Direction == "down" {
			e.CurrentFloor--
		}
		time.Sleep(5 * time.Second)
	}
	e.Moving = false
}

func (e *Elevator) OpenDoors() {
	e.DoorsOpen = true
	time.Sleep(2 * time.Second)
	e.DoorsOpen = false
}
