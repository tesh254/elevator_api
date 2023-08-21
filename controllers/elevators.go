package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"elevator/queries"
	"elevator/ws"
)

type RequestData struct {
	ToFloor    int `json:"to_floor"`
	ElevatorID int `json:"elevator_id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func ElevatorHandler(w http.ResponseWriter, req *http.Request, queries queries.Queries) {
	var requestData RequestData

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	elevatorID := requestData.ElevatorID
	toFloor := requestData.ToFloor

	elevator, err := queries.GetElevatorById(elevatorID)
	if err != nil {
		errorResponse := ErrorResponse{Message: fmt.Sprintf("Elevator was not found: %s", err)}
		errorJSON, _ := json.Marshal(errorResponse)
		http.Error(w, string(errorJSON), http.StatusNotFound)
		return
	}

	go ws.SendData(fmt.Sprintf(`{"id": %d, "current_floor": %d, "direction": "%s", "moving": %t, "doors_open": %t, "building_id": %d, "to_floor": %d }`, elevator.ID, elevator.Floor, elevator.Direction, elevator.Moving, elevator.IsDoorsOpen, elevator.BuildingID, toFloor))

	response := map[string]string{"message": "Elevator called"}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
