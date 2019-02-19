package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
  "fmt"
)

type daysReadResponse struct {
  Days []*Day `json:"days"`
}

func daysHandler(r *mux.Router) {
  r.HandleFunc("/days/{year:[0-9]+}/{month:[0-9]+}", daysIndexHandler).Methods("GET")
  r.HandleFunc("/days/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}", daysUpdateHandler).Methods("PATCH")
}

func daysIndexHandler(w http.ResponseWriter, r *http.Request) {
  // Parse token info
  _, err := getAuthTokenClaims(r)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    sendJson(w, errorResponse{ Error: "Invalid token" })
    log.Printf("Invalid token: %v", err)
    return
  }

  vars := mux.Vars(r)

  // Get days
  var days []*Day
  query := "EXTRACT(YEAR FROM date) = ? and EXTRACT(MONTH FROM date) = ?"
  db.Where(query, vars["year"], vars["month"]).Find(&days)

  // Respond
  response := daysReadResponse{ Days: days }
  w.WriteHeader(http.StatusOK)
  sendJson(w, response)
}

func daysUpdateHandler(w http.ResponseWriter, r *http.Request) {
  // Parse token info
  _, err := getAuthTokenClaims(r)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    sendJson(w, errorResponse{ Error: "Invalid token" })
    log.Printf("Invalid token: %v", err)
    return
  }

  vars := mux.Vars(r)

  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body map[string]*int
  err = decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  // Find day
  date := fmt.Sprintf(
    "%s/%s/%s",
    vars["month"],
    vars["day"],
    vars["year"])
  var day = &Day{}
  res := db.Where("date = ?", date).First(&day)
  if res.RecordNotFound() {
    // Create the day if it doesn't exists
    db.Create(&day)
  }
  if res.Error != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error querying db" })
    log.Printf("Error querying db: %v", err)
    return
  }

  // Update values
  if val, ok := body["bmr"]; ok {
    if val == nil {
      day.Bmr = nil
    } else {
      day.Bmr = val
    }
  }
  if val, ok := body["caloriesIn"]; ok {
    if val == nil {
      day.CaloriesIn = nil
    } else {
      day.CaloriesIn = val
    }
  }
  if val, ok := body["caloriesOut"]; ok {
    if val == nil {
      day.CaloriesOut = nil
    } else {
      day.CaloriesOut = val
    }
  }
  if val, ok := body["caloriesGoal"]; ok {
    if val == nil {
      day.CaloriesGoal = nil
    } else {
      day.CaloriesGoal = val
    }
  }
  if val, ok := body["milesRun"]; ok {
    if val == nil {
      day.MilesRun = nil
    } else {
      day.MilesRun = val
    }
  }
  if val, ok := body["milesRunGoal"]; ok {
    if val == nil {
      day.MilesRunGoal = nil
    } else {
      day.MilesRunGoal = val
    }
  }
  if val, ok := body["drinks"]; ok {
    if val == nil {
      day.Drinks = nil
    } else {
      day.Drinks = val
    }
  }
  if val, ok := body["drinksGoal"]; ok {
    if val == nil {
      day.DrinksGoal = nil
    } else {
      day.DrinksGoal = val
    }
  }

  // Update record
  if err := db.Save(&day).Error; err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error updating db" })
    log.Printf("Error updating db: %v", err)
    return
  }

  // Respond
  w.WriteHeader(http.StatusOK)
}
