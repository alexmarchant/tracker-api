package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
)

func usersHandler(r *mux.Router) {
  r.HandleFunc("/users", usersCreateHandler).Methods("POST")
}

type usersCreateRequest struct {
  Email string `json:"email"`
  Password string `json:"password"`
  PasswordConfirmation string `json:"passwordConfirmation"`
}

type usersCreateResponse struct {
  Token string `json:"token"`
}

func usersCreateHandler(w http.ResponseWriter, r *http.Request) {
  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body usersCreateRequest
  err := decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  // Validate request
  if body.Email == "" || body.Password == "" || body.PasswordConfirmation == "" {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Missing required params" })
    log.Print("Missing required params")
    return
  }

  // Check password matched passwordConfirmation
  if body.Password != body.PasswordConfirmation {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Password doesn't match confirmation" })
    log.Print("Password doesn't match confirmation")
    return
  }

  // Check password length
  if len(body.Password) < 6 {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "Password must be at least 6 characters long" })
    log.Print("Password must be at least 6 characters long")
    return
  }

  // Check if user exists
  var id int64
  err = db.QueryRow("SELECT id FROM users WHERE email = $1", body.Email).Scan(&id)
  if err == nil {
    w.WriteHeader(http.StatusBadRequest)
    sendJson(w, errorResponse{ Error: "User already exists" })
    log.Print("User already exists")
    return
  }

  // Hash pw
  passwordHash := hashAndSalt(body.Password)

  // Create user
  err = db.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", body.Email, passwordHash).Scan(&id)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    sendJson(w, errorResponse{ Error: "Error writing user data to database" })
    log.Printf("Error writing user data to database: %v", err)
    return
  }

  // Create JWT Token
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
  sendJson(w, usersCreateResponse{ Token: token })
}
