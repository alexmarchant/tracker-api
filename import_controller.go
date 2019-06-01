package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func importHandler(r *mux.Router) {
	r.HandleFunc("/import", importCreateHandler).Methods("POST")
}

func importCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	decoder := json.NewDecoder(r.Body)
	var body map[string]float64
	err := decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendJson(w, errorResponse{Error: "Error parsing request"})
		log.Print("Error parsing request")
		return
	}

	// Map request to time and add to map
	// years has this structure:
	// years = { 2019: { 4: { 23: [123, 143] } } }
	dateWeights := map[string][]float64{}
	for dateString, weight := range body {
		// Dates should be in this format "2014-11-12T11:45:26.371Z"
		reqTime, err := time.Parse(time.RFC3339, dateString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			sendJson(w, errorResponse{Error: "Invalid date"})
			return
		}

		// Convert utc to est
		location, err := time.LoadLocation("EST")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			sendJson(w, errorResponse{Error: "Internal server error"})
			return
		}
		estTime := reqTime.In(location)

		// Get day ints
		year, monthVal, day := estTime.Date()
		month := int(monthVal)

		// Ensure day data
		dateString := fmt.Sprintf("%d/%d/%d", month, day, year)
		if _, ok := dateWeights[dateString]; !ok {
			dateWeights[dateString] = []float64{}
		}

		// Add weight
		dateWeights[dateString] = append(dateWeights[dateString], weight)
	}

	// Dedupe values and pick the min weight of any given day
	dateWeightsSet := map[string]float64{}
	for dateString, weights := range dateWeights {
		// Find min
		min := weights[0]
		for _, weight := range weights {
			if weight < min {
				min = weight
			}
		}

		// Set weight
		dateWeightsSet[dateString] = min
	}

	// Get relevant days
	queryDates := []string{}
	for dateString := range dateWeightsSet {
		queryDates = append(queryDates, dateString)
	}
	dbDays := []Day{}
	res := db.Where("date IN (?)", queryDates).Find(&dbDays)
	if res.Error != nil && !res.RecordNotFound() {
		w.WriteHeader(http.StatusInternalServerError)
		sendJson(w, errorResponse{Error: "Error querying db"})
		log.Printf("Error querying db: %v", res.Error)
		return
	}

	for dateString, weight := range dateWeightsSet {
		// Already exists?
		exists := false
		for _, dbDay := range dbDays {
			// Parse time
			dbTime, err := time.Parse(time.RFC3339, dbDay.Date)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				sendJson(w, errorResponse{Error: "Internal server error"})
				return
			}
			year, monthVal, day := dbTime.Date()
			month := int(monthVal)
			dbDateString := fmt.Sprintf("%d/%d/%d", month, day, year)

			// Check if exists
			if dbDateString == dateString {
				// Check if unchanged, skip if so
				if *dbDay.Weight == weight {
					exists = true
					continue
				}

				// Update row
				dbDay.Weight = &weight
				if err := db.Save(&dbDay).Error; err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					sendJson(w, errorResponse{Error: "Error updating db"})
					log.Printf("Error updating db: %v", err)
					return
				}
				exists = true
				continue
			}
		}

		// Check if we found it in previous loop
		if exists {
			continue
		}

		// Create new
		newDay := Day{
			Date:   dateString,
			Weight: &weight,
		}
		if err := db.Create(&newDay).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			sendJson(w, errorResponse{Error: "Error updating db"})
			log.Printf("Error updating db: %v", err)
			return
		}
	}

	// Respond
	w.WriteHeader(http.StatusOK)
}
