package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
  "fmt"
)

type day struct {
  Date string `json:"date"`
  Bmr int `json:"bmr"`
  CaloriesIn *int `json:"caloriesIn,omitempty"`
  CaloriesOut *int `json:"caloriesOut,omitempty"`
}

type daysReadResponse struct {
  Days []*day `json:"days"`
}

type daysCreateRequest struct {
  Bmr *int `json:"bmr,omitempty"`
  CaloriesIn *int `json:"caloriesIn,omitempty"`
  CaloriesOut *int `json:"caloriesOut,omitempty"`
}

func daysHandler(r *mux.Router) {
  r.HandleFunc("/days/{year:[0-9]+}/{month:[0-9]+}", daysIndexHandler).Methods("GET")
  r.HandleFunc("/days/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}", daysCreateHandler).Methods("POST")
}

func daysIndexHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  // Get calories
  query := `
    SELECT TO_CHAR(date, 'MM/DD/YYYY'), bmr, calories_in, calories_out
    FROM days
    WHERE EXTRACT(YEAR FROM date) = $1 and EXTRACT(MONTH FROM date) = $2
  `
  rows, err := db.Query(
    query,
    vars["year"],
    vars["month"])
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  var days = []*day{}

  // Iterate over messages
  for rows.Next() {
    var aDay day

    err := rows.Scan(&aDay.Date, &aDay.Bmr, &aDay.CaloriesIn, &aDay.CaloriesOut)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      sendJson(w, errorResponse{ error: "Error querying database" })
      log.Printf("Error querying database: %v", err)
      return
    }

    days = append(days, &aDay)
  }

  // Check iteration for errors
  if err := rows.Err(); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  // Respond
  response := daysReadResponse{ Days: days }
  w.WriteHeader(http.StatusOK)
  sendJson(w, response)
}

func daysCreateHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body daysCreateRequest
  err := decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  query := `
    INSERT INTO days (date, bmr, calories_in, calories_out) 
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (date) DO UPDATE 
      SET date = excluded.date, 
        bmr = excluded.bmr,
        calories_in = excluded.calories_in,
        calories_out = excluded.calories_out;
  `
  date := fmt.Sprintf("%s/%s/%s", vars["month"], vars["day"], vars["year"])
  _, err = db.Exec(query, date, &body.Bmr, &body.CaloriesIn, &body.CaloriesOut)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ error: "Error updating database" })
    log.Printf("Error updating picks: %v", err)
    return
  }

  w.WriteHeader(http.StatusCreated)
}