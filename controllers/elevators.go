package controllers

import (
	"encoding/json"
	"net/http"
)

type RequestData struct {
	ToFloor    int `json:"to_floor"`
	ElevatorID int `json:"elevator_id"`
}

func ElevatorHandler(w http.ResponseWriter, req *http.Request) {
	var requestData RequestData

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// toFloor := requestData.ToFloor
	// elevatorID := requestData.ElevatorID
}
