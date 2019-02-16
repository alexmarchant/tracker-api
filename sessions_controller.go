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

type sessionsCreateRequest struct {
  Email string `json:"email"`
  Password string `json:"password"`
}

type sessionsCreateResponse struct {
  Token string `json:"token"`
}

func sessionsCreateHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request
  decoder := json.NewDecoder(r.Body)
  var body sessionsCreateRequest
  err := decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  // Validate request
  if body.Email == "" || body.Password == "" {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Missing required params" })
    log.Print("Missing required params")
    return
  }

  // Get user
  var id int64
  var passwordHash string
  err = db.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", body.Email).Scan(&id, &passwordHash)
  if err != nil {
    if err.Error() == "sql: no rows in result set" {
      w.WriteHeader(http.StatusBadRequest)
      sendJson(w, errorResponse{ Error: "No user found for that email" })
      log.Print("No user found for that email")
    } else {
      w.WriteHeader(http.StatusInternalServerError)
      sendJson(w, errorResponse{ Error: "Error finding user" })
      log.Printf("Error finding user %v", err)
    }
    return
  }

  // Check password
  if !comparePasswords(body.Password, passwordHash) {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Wrong password" })
    log.Print("Wrong password")
    return
  }

  // Create token
  claims := &tokenClaims{
    Id: id,
    Email: body.Email,
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
