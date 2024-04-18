package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func logRequestMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method: %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func isTokenProvided(r *http.Request) bool {
	token := r.Header.Get("Authorization")

	return token != ""
}

func getUserFromRequest(r *http.Request) (UserModel, error) {
	user := UserModel{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		return UserModel{}, fmt.Errorf("Error decoding user: %w", err)
	}

	return user, nil
}

func createHandler(controller Controller) http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("POST /user/create", handleUserCreate(controller))
	handler.HandleFunc("GET /user/habits", handleUserGet(controller))

	handler.HandleFunc("PUT /user/habits", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("updating user habits"))
	})

	handler.HandleFunc("PUT /user/token", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("revoking user token"))
	})

	return logRequestMiddleware(handler)
}

func handleUserCreate(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, "Failed to decode or missing data.", http.StatusInternalServerError)

			return
		}

		bytes, err := json.Marshal(user)
		if err != nil {
			log.Printf("error when marshaling user: %v", err)

			return
		}

		exists, err := controller.CreateNewUser(user)
		if err != nil {
			http.Error(w, "Failed to create user.", http.StatusInternalServerError)

			return
		}

		if exists {
			w.WriteHeader(http.StatusOK)

			return
		}

		if _, err := w.Write(bytes); err != nil {
			log.Printf("Failed to write bytes error: %v", err)
		}
	}
}

func handleUserGet(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isTokenProvided(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, "Failed to decode or missing data.", http.StatusInternalServerError)

			return
		}

		habits, err := controller.GetUserHabits(user)
		if errors.Is(err, ErrInccorectInput) {
			// TODO: handle some error code and message
			return
		}

		bytes, _ := json.Marshal(habits)

		_, _ = w.Write(bytes)
	}
}
