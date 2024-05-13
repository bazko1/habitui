package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const missingUserInputErrMessage = "Failed to decode user or missing data."

var ErrBadAuthorizationHeader = errors.New("bad authorization header")

func logRequestMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method: %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func getBearerToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")

	if len(splitToken) != 2 {
		return "", ErrBadAuthorizationHeader
	}

	return strings.TrimSpace(splitToken[1]), nil
}

func getUserFromRequest(r *http.Request) (UserModel, error) {
	user := UserModel{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		return UserModel{}, fmt.Errorf("error decoding user: %w", err)
	}

	return user, nil
}

func createHandler(controller Controller) http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("POST /user/create", handlePostUserCreate(controller))
	handler.HandleFunc("POST /user/login", handlePostUserLogin(controller))
	handler.HandleFunc("GET /user/habits", handleGetUserHabits(controller))
	handler.HandleFunc("PUT /user/habits", handlePutUserHabits(controller))

	// TODO: for now use single one time token auth that cannot be revoked
	// implement more sophisticated authentication later
	// handler.HandleFunc("PUT /user/token", handlePutUserToken(controller))

	return logRequestMiddleware(handler)
}

func handlePostUserCreate(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO this should get username, email, password
		// and create new user and create jwt token
		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, missingUserInputErrMessage, http.StatusInternalServerError)

			return
		}

		newUser, err := controller.CreateNewUser(user)
		if errors.Is(err, ErrInccorectInput) {
			http.Error(w, fmt.Sprintf("Incorrect input error: %v", err), http.StatusUnprocessableEntity)

			return
		}

		if errors.Is(err, ErrUsernameExists) {
			w.WriteHeader(http.StatusNoContent)
			_, _ = w.Write([]byte("{}"))

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

		w.WriteHeader(http.StatusCreated)

		if _, err := w.Write(bytes); err != nil {
			log.Printf("Failed to write bytes error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)

			return
		}
	}
}

func handleGetUserHabits(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: create JWT token validation
		// token, err := getBearerToken(r)

		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, missingUserInputErrMessage, http.StatusInternalServerError)

			return
		}

		habits, err := controller.GetUserHabits(user)
		if errors.Is(err, ErrNonExistentUserOrPassword) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

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

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bytes)
	}
}

func handlePutUserHabits(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: this does not work properly need to refactor so that
		// this works with json from habit module otherwise the
		// tasks will not be properly filled with data and will be
		// just nil initialized
		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, missingUserInputErrMessage, http.StatusInternalServerError)

			return
		}

		err = controller.UpdateUserHabits(user, user.Habits)
		if errors.Is(err, ErrNonExistentUserOrPassword) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		if err != nil {
			log.Printf("Updating user habits error: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handlePostUserLogin(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRequest(r)
		if err != nil {
			log.Printf("Error getting user from request: %v", err)
			http.Error(w, missingUserInputErrMessage, http.StatusInternalServerError)

			return
		}

		if !controller.IsValid(user) {
			http.Error(w, ErrNonExistentUserOrPassword.Error(), http.StatusUnauthorized)

			return
		}

		token, err := generateJWT(user.Username)
		if err != nil {
			log.Printf("error generating jwt: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		// TODO: create single response template that has versioning
		// status fields and some data fields, because this will be
		// maintenance and support mess.
		bytes, err := json.Marshal(map[string]string{"token": token})
		if err != nil {
			log.Printf("error when marshaling user: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(bytes)
	}
}
