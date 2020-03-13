package app

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			message := map[string]interface{}{"error": "Missing Authorization Header"}
			messageJSON, _ := json.Marshal(message)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(messageJSON)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := verifyToken(tokenString)
		if err != nil {
			message := map[string]interface{}{"error": "Error verifying JWT token: " + err.Error()}
			messageJSON, _ := json.Marshal(message)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(messageJSON)
			return
		}
		username := claims.(jwt.MapClaims)["username"].(string)

		r.Header.Set("username", username)

		next.ServeHTTP(w, r)
	})
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

// Get : Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post : Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put : Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete : Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// GetWithAuth : Wrap the router for GET method
func (a *App) GetWithAuth(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.Handle(path, authMiddleware(http.HandlerFunc(f))).Methods("GET")
}

// PostWithAuth : Wrap the router for POST method
func (a *App) PostWithAuth(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.Handle(path, authMiddleware(http.HandlerFunc(f))).Methods("POST")
}

// PutWithAuth : Wrap the router for PUT method
func (a *App) PutWithAuth(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.Handle(path, authMiddleware(http.HandlerFunc(f))).Methods("PUT")
}

// DeleteWithAuth : Wrap the router for DELETE method
func (a *App) DeleteWithAuth(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.Handle(path, authMiddleware(http.HandlerFunc(f))).Methods("DELETE")
}

// Static handle router for static file
func (a *App) Static() {
	a.Router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("uploads/"))))
}
