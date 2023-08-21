package ws

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func SendData(message string) {
	serverAddr := os.Getenv("WS_URL")

	baseURL, err := url.Parse(serverAddr)

	if err != nil {
		log.Fatal("URL parsing error:", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(baseURL.String(), nil)
	if err != nil {
		log.Fatal("WebSocket connection error:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatal("WebSocket write error:", err)
	}

	for {
		_, receivedMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", receivedMsg)

		time.Sleep(time.Second)
	}
}
