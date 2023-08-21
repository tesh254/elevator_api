package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"elevator/core"
	"elevator/db"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*websocket.Conn]bool)
)

func WsHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			delete(clients, conn)
			break
		}

		var elevatorUpdates core.Elevator

		err = json.Unmarshal([]byte(msg), &elevatorUpdates)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		elevatorUpdates.WsConn = conn
		elevatorUpdates.DbConn = database
		elevatorUpdates.Clients = clients

		elevatorUpdates.Start()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}

func StartWebSocketServer() {
	port := 8081
	database := db.ConnectToDatabase(false)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		WsHandler(w, r, database)
	})

	log.Printf("🏃‍♂️💨 :::WebsoketServer is starting on port %d:::\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("WebSocket server error: ", err)
	}
}
