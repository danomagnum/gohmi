package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"flag"
	"net/http"
)

var basicHTTPAuthEnabled = flag.Bool("basicHTTPAuth", false, "Enable basic HTTP Authentication")
var basicUser = flag.String("basicHTTPUser", "admin", "username for basic http auth")
var basicPass = flag.String("basicHTTPPass", "password", "password for basic http auth")

// middleware to enfore User and Pass through basic HTTP authentication
type basicHTTPAuth struct {
	User string
	Pass string

	Next http.Handler
}

func (a *basicHTTPAuth) Handler(next http.Handler) http.Handler {
	a.Next = next
	return a
}

func (a *basicHTTPAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		expectedUsernameHash := sha256.Sum256([]byte(a.User))
		expectedPasswordHash := sha256.Sum256([]byte(a.Pass))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {
			a.Next.ServeHTTP(w, r)
			return
		}
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
