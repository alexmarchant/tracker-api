package main

import (
	"database/sql"
	_ "github.com/lib/pq"
  "os"
  "log"
)

var db *sql.DB

func connectDB() {
  pgURL := os.Getenv("PG_URL")
  if pgURL == "" {
    log.Fatal("Missing PG_URL")
  }

  // Create connection pointer
  var err error
  db, err = sql.Open("postgres", pgURL)
  if err != nil {
    log.Fatal(err)
  }

  // Actually establish connection
  err = db.Ping()
  if err != nil {
    log.Print("Error connecting to database:")
    log.Fatal(err)
  }
}
