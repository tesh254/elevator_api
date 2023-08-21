package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"

	"elevator/db"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	godotenv.Load()

	database := db.ConnectToDatabase()

	defer database.Close()

	router := http.NewServeMux()

	router.HandleFunc("/hello", hello)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // You can specify specific origins instead of "*"
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := 8090

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: corsHandler(router),
	}

	log.Printf("🏃‍♂️ :::Server is starting on port %d:::\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
