package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "strings"
  "log"
  "fmt"
)

type day struct {
  Date string `json:"date"`
  Bmr *int `json:"bmr,omitempty"`
  CaloriesIn *int `json:"caloriesIn,omitempty"`
  CaloriesOut *int `json:"caloriesOut,omitempty"`
  MilesRun *int `json:"milesRun,omitempty"`
}

type daysReadResponse struct {
  Days []*day `json:"days"`
}

type daysUpdateRequest struct {
  Bmr *int `json:"bmr,omitempty"`
  CaloriesIn *int `json:"caloriesIn,omitempty"`
  CaloriesOut *int `json:"caloriesOut,omitempty"`
  MilesRun *int `json:"milesRun,omitempty"`
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

  // Get calories
  query := `
    SELECT TO_CHAR(date, 'MM/DD/YYYY'), bmr, calories_in, calories_out, miles_run
    FROM days
    WHERE EXTRACT(YEAR FROM date) = $1 and EXTRACT(MONTH FROM date) = $2
  `
  rows, err := db.Query(
    query,
    vars["year"],
    vars["month"])
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  var days = []*day{}

  // Iterate over messages
  for rows.Next() {
    var aDay day

    err := rows.Scan(&aDay.Date, &aDay.Bmr, &aDay.CaloriesIn, &aDay.CaloriesOut, &aDay.MilesRun)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      sendJson(w, errorResponse{ Error: "Error querying database" })
      log.Printf("Error querying database: %v", err)
      return
    }

    days = append(days, &aDay)
  }

  // Check iteration for errors
  if err := rows.Err(); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

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
  var body daysUpdateRequest
  err = decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  } // Get day
  date := fmt.Sprintf("%s/%s/%s", vars["month"], vars["day"], vars["year"])
  var aDay day
  query := `
    SELECT date, bmr, calories_in, calories_out, miles_run
    FROM days
    WHERE date = $1
  `
  err = db.QueryRow(query, date).Scan(&aDay.Date, &aDay.Bmr, &aDay.CaloriesIn, &aDay.CaloriesOut, &aDay.MilesRun)
  if err != nil {
    if err.Error() == "sql: no rows in result set" {
      createDay(date, body, w)
    } else {
      w.WriteHeader(http.StatusInternalServerError)
      sendJson(w, errorResponse{ Error: "Error querying database" })
      log.Printf("Error querying database: %v", err)
    }
    return
  }

  updateDay(date, aDay, body, w)
}

func createDay(date string, body daysUpdateRequest, w http.ResponseWriter) {
  query := `
    INSERT INTO days (date, bmr, calories_in, calories_out, miles_run) 
    VALUES ($1, $2, $3, $4, $5)
  `
  _, err := db.Exec(query, date, &body.Bmr, &body.CaloriesIn, &body.CaloriesOut, &body.MilesRun)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error creating day" })
    log.Printf("Error creating day: %v", err)
    return
  }

  w.WriteHeader(http.StatusCreated)
}

func updateDay(date string, aDay day, body daysUpdateRequest, w http.ResponseWriter) {
  query := "UPDATE days SET "
  var querySet []string
  var queryParams []interface{}

  if body.Bmr != nil {
    queryParams = append(queryParams, body.Bmr)
    querySet = append(querySet, fmt.Sprintf("bmr = $%d", len(queryParams)))
  }
  if body.CaloriesIn != nil {
    queryParams = append(queryParams, body.CaloriesIn)
    querySet = append(querySet, fmt.Sprintf("calories_in = $%d", len(queryParams)))
  }
  if body.CaloriesOut != nil {
    queryParams = append(queryParams, body.CaloriesOut)
    querySet = append(querySet, fmt.Sprintf("calories_out = $%d", len(queryParams)))
  }
  if body.MilesRun != nil {
    queryParams = append(queryParams, body.MilesRun)
    querySet = append(querySet, fmt.Sprintf("miles_run = $%d", len(queryParams)))
  }

  if len(queryParams) == 0 {
    w.WriteHeader(http.StatusOK)
    return
  }

  query += strings.Join(querySet, ", ")
  queryParams = append(queryParams, date)
  query += fmt.Sprintf(" WHERE date = $%d", len(queryParams))

  _, err := db.Exec(query, queryParams...)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error updating day" })
    log.Printf("Error updating day: %v", err)
    return
  }

  w.WriteHeader(http.StatusOK)
}
