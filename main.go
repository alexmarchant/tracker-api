package main

import (
  "github.com/rs/cors"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "net/http"
  "os"
)

func main() {
  r := mux.NewRouter()

  // DB
  connectDB()

  // Routes
  daysHandler(r)

  // CORS
  cors := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders: []string{"*"},
  })
  rCors := cors.Handler(r)

  // Loggging
  rLogging := handlers.LoggingHandler(os.Stdout, rCors)

  // Start server
  http.ListenAndServe(":3000", rLogging)
}
