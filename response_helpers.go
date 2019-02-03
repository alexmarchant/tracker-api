package main

import (
  "net/http"
  "encoding/json"
)

type errorResponse struct {
  error string `json:"error"`
}

func sendJson(w http.ResponseWriter, data interface{}) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(data)
}
