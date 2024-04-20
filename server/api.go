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

	handler.HandleFunc("POST /user/create", handlePostUserCreate(controller))
	handler.HandleFunc("GET /user/habits", handleGetUserHabits(controller))

	handler.HandleFunc("PUT /user/habits", handlePutUserHabits(controller))

	handler.HandleFunc("PUT /user/token", handlePutUserToken(controller))

	return logRequestMiddleware(handler)
}

func handlePostUserCreate(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, "Failed to decode or missing data.", http.StatusInternalServerError)

			return
		}

		newUser, err := controller.CreateNewUser(user)

		if errors.Is(err, ErrUsernameExists) {
			w.Write([]byte("{}"))
			w.WriteHeader(http.StatusNoContent)

			return
		}

		if err != nil {
			log.Printf("error when creating user: %v", err)
			http.Error(w, "Failed to create user.", http.StatusInternalServerError)

			return
		}

		bytes, err := json.Marshal(newUser)
		if err != nil {
			log.Printf("error when marshaling user: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		if _, err := w.Write(bytes); err != nil {
			log.Printf("Failed to write bytes error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleGetUserHabits(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, "Failed to decode or missing data.", http.StatusInternalServerError)

			return
		}

		habits, err := controller.GetUserHabits(user)
		if errors.Is(err, ErrInccorectInput) {
			http.Error(w, fmt.Sprintf("Incorrect input error: %v", err), http.StatusUnprocessableEntity)

			return
		}

		if err != nil {
			log.Printf("Getting user habits error: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		bytes, err := json.Marshal(habits)
		if err != nil {
			log.Printf("error when marshaling user: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(bytes)
	}
}

func handlePutUserHabits(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, "Failed to decode or missing data.", http.StatusInternalServerError)

			return
		}

		habits, err := controller.GetUserHabits(user)
		if errors.Is(err, ErrInccorectInput) {
			http.Error(w, fmt.Sprintf("Incorrect input error: %v", err), http.StatusUnprocessableEntity)

			return
		}

		if err != nil {
			log.Printf("Getting user habits error: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		bytes, err := json.Marshal(habits)
		if err != nil {
			log.Printf("error when marshaling user: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(bytes)
	}
}

func handlePutUserToken(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, "Failed to decode or missing data.", http.StatusInternalServerError)

			return
		}

		habits, err := controller.GetUserHabits(user)
		if errors.Is(err, ErrInccorectInput) {
			http.Error(w, fmt.Sprintf("Incorrect input error: %v", err), http.StatusUnprocessableEntity)

			return
		}

		if err != nil {
			log.Printf("Getting user habits error: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		bytes, err := json.Marshal(habits)
		if err != nil {
			log.Printf("error when marshaling user: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(bytes)
	}
}
