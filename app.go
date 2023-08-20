package main

import (
	"fmt"
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

	http.ListenAndServe(":8090", corsHandler((router)))
}
