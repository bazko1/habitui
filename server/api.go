package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var controller = NewInMemoryController()

func logRequestMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method: %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func isTokenProvided(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		return false
	}
	return true
}

func getUserFromRequest(r *http.Request) (UserModel, error) {
	user := UserModel{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		return UserModel{}, fmt.Errorf("Error decoding user: %w", err)
	}
	return user, nil
}

func createHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("POST /user/create", handleUserCreate)
	handler.HandleFunc("GET /user/habits", handleUserGet)

	handler.HandleFunc("PUT /user/habits", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("updating user habits"))
	})

	handler.HandleFunc("PUT /user/token", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("revoking user token"))
	})
	return logRequestMiddleware(handler)
}

func handleUserCreate(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		log.Printf("Error getting user from request: %v", err)
		http.Error(w, "Failed to decode or missing data", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		log.Printf("error when marshaling user: %v", err)
		return
	}

	controller.CreateNewuser(user.Username, user.Email)
	_, _ = w.Write(bytes)
}

func handleUserGet(w http.ResponseWriter, r *http.Request) {
	if !isTokenProvided(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
