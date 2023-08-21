package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"elevator/queries"
)

type LogRequestData struct {
	ElevatorID int `json:"elevator_id"`
}

type LogErrorResponse struct {
	Message string `json:"message"`
}

func LogsHandler(w http.ResponseWriter, req *http.Request, queries queries.Queries) {
	var requestData LogRequestData

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	elevatorID := requestData.ElevatorID

	logs, err := queries.GetLogsByElevator(elevatorID)
	if err != nil {
		errorResponse := LogErrorResponse{Message: fmt.Sprintf("Elevator was not found: %s", err)}
		errorJSON, _ := json.Marshal(errorResponse)
		http.Error(w, string(errorJSON), http.StatusNotFound)
		return
	}

	response := logs
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
