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

func getBearerToken(r *http.Request) (map[string]any, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")

	//// nolint:gomnd // if string has 'Bearer' keyword and token string
	// result of split should have empty string and token.
	if len(splitToken) != 2 {
		return map[string]any{}, ErrBadAuthorizationHeader
	}

	s := strings.TrimSpace(splitToken[1])

	claims, err := parseAndValidateJWT(s)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error getting bearer token: %w", err)
	}

	return claims, err
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

		_, err = controller.CreateNewUser(user)
		if errors.Is(err, ErrInccorectInput) {
			http.Error(w, fmt.Sprintf("Incorrect input error: %v", err), http.StatusUnprocessableEntity)

			return
		}

		if errors.Is(err, ErrUsernameExists) {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		if err != nil {
			log.Printf("error when creating user: %v", err)
			http.Error(w, "Failed to create user.", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusCreated)

		return
	}
}

func handleGetUserHabits(controller Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getBearerToken(r)
		if err != nil {
			log.Printf("err when getting bearer token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		user, ok := controller.GetUserByName(claims["username"].(string))
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

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
		claims, err := getBearerToken(r)
		if err != nil {
			log.Printf("err when getting bearer token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		user, ok := controller.GetUserByName(claims["username"].(string))
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

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

		tokenMap, err := generateJWT(user.Username)
		if err != nil {
			log.Printf("error generating jwt: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		bytes, err := json.Marshal(tokenMap)
		if err != nil {
			log.Printf("error when marshaling user: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(bytes)
	}
}
