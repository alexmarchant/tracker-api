package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "os"
  "log"
)

var db *gorm.DB

func connectDB() {
  pgURL := os.Getenv("PG_URL")
  if pgURL == "" {
    log.Fatal("Missing PG_URL")
  }

  // Connect
  var err error
  db, err = gorm.Open("postgres", pgURL)
  if err != nil {
    panic("failed to connect database")
  }

  // Config
  db.LogMode(true)

  // Migrate the schema
  db.AutoMigrate(&Day{})
  db.AutoMigrate(&User{})
}
