package main

import (
  "github.com/dgrijalva/jwt-go"
  "golang.org/x/crypto/bcrypt"
  "net/http"
  "strings"
  "errors"
  "fmt"
  "log"
	"os"
)

var tokenSecret string

type tokenClaims struct {
  Id int64 `json:"id"`
  Email string `json:"email"`
  jwt.StandardClaims
}

func getTokenSecret() {
	// Ensure secret
	tokenSecret = os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("Missing TOKEN_SECRET")
	}
}

func hashAndSalt(password string) string {
  passwordBytes := []byte(password)
  hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
  if err != nil {
    log.Println(err)
  }
  return string(hash)
}

func makeToken(claims *tokenClaims) (string, error) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString([]byte(tokenSecret))
}

func comparePasswords(password string, passwordHash string) bool {
  passwordBytes := []byte(password)
  hashBytes := []byte(passwordHash)
  err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
  if err != nil {
      log.Println(err)
      return false
  }

  return true
}

func getAuthTokenString(r *http.Request) (string, error) {
  var tokenString string

  // Get header
  val := r.Header["Authorization"]
  if len(val) == 0 {
    return tokenString, errors.New("Header not found")
  }

  // Parse out bearer
  tokenString = strings.TrimPrefix(val[0], "Bearer ")
  return tokenString, nil
}

func getAuthTokenClaimsFromString(tokenString string) (*tokenClaims, error) {
  // Decode token
  token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
    // Validate method is what we expect
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(tokenSecret), nil
  })
  if err != nil {
    return nil, err
  }

  // Get claims
  if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
    return claims, nil
  } else {
    return nil, err
  }
}

func getAuthTokenClaims(r *http.Request) (*tokenClaims, error) {
  // Get token string
  tokenString, err := getAuthTokenString(r)
  if err != nil {
    return nil, err
  }

  return getAuthTokenClaimsFromString(tokenString)
}
