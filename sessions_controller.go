package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
)

func sessionsHandler(r *mux.Router) {
  r.HandleFunc("/sessions", sessionsCreateHandler).Methods("POST")
}

func sessionsCreateHandler(w http.ResponseWriter, r *http.Request) {
  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body map[string]interface{}
  err := decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  // Validate request
  email, ok := body["email"]
  if !ok || email == "" {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Missing required params" })
    log.Print("Missing required params")
    return
  }
  password, ok := body["password"]
  if !ok || password == "" {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Missing required params" })
    log.Print("Missing required params")
    return
  }

  // Get user
  var user User
  if db.Where("email = ?", email).First(&user).RecordNotFound() {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "No user found for that email" })
    log.Print("No user found for that email")
    return
  }

  // Check password
  if !comparePasswords(password.(string), user.PasswordHash) {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Wrong password" })
    log.Print("Wrong password")
    return
  }

  // Create token
  type sessionsCreateResponse struct {
    Token string `json:"token"`
  }
  claims := &tokenClaims{
    Id: user.ID,
    Email: user.Email,
  }
  token, err := makeToken(claims)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error creating token" })
    log.Printf("Error creating token: %s", err)
    return
  }

  // Respond
  w.WriteHeader(http.StatusCreated)
  sendJson(w, sessionsCreateResponse{ Token: token })
}
