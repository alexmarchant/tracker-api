package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	_ = godotenv.Load()
	r := mux.NewRouter()

	// DB
	connectDB()

	// Setup auth
	getTokenSecret()

	// Routes
	daysHandler(r)
	sessionsHandler(r)
	usersHandler(r)
	importHandler(r)

	// CORS
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"*"},
	})
	rCors := cors.Handler(r)

	// Loggging
	rLogging := handlers.LoggingHandler(os.Stdout, rCors)

	// Start server
	port := 3000
	fmt.Printf("Server running on port %d\n", port)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), rLogging)
}
