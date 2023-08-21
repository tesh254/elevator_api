package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"

	"elevator/controllers"
	"elevator/db"
	"elevator/queries"
	"elevator/ws"
)

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		next(w, req)

		duration := time.Since(start)
		fmt.Printf("[%s] %s %s %v\n", time.Now().Format("2006-01-02 15:04:05"), req.Method, req.URL.Path, duration)
	}
}

func main() {
	godotenv.Load()

	database := db.ConnectToDatabase(true)

	var genSeed bool
	var run bool
	var dropTables bool

	flag.BoolVar(&genSeed, "gen-seed", false, "Generate seed data")
	flag.BoolVar(&run, "run", false, "Run server and WebSocket")
	flag.BoolVar(&dropTables, "drop-tables", false, "Drop all tables")

	flag.Parse()

	if dropTables {
		db.DropTables(database)
		return
	}

	if genSeed {
		db.SeedTable(database)
		return
	}

	if run {
		router := http.NewServeMux()

		var queries queries.Queries

		queries.Database = database

		router.HandleFunc("/call", logRequest(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodPost {
				controllers.ElevatorHandler(w, req, queries)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))

		router.HandleFunc("/logs", logRequest(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodGet {
				controllers.LogsHandler(w, req, queries)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))

		corsHandler := handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		)

		go ws.StartWebSocketServer()

		port := 8090

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: corsHandler(router),
		}

		defer database.Close()

		log.Printf("🏃‍♂️ :::Server is starting on port %d:::\n", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Server error: ", err)
		}
	}
}
