package main

// APIError field error message
import (
	"encoding/json"
	"net/http"
)

// APIError message for building API error
type APIError struct {
	Field   string
	Message string
}

func apiErrorHandler(w http.ResponseWriter, r *http.Request, status int, errors []APIError) {
	w.WriteHeader(status)

	errorBody := map[string]interface{}{}
	for _, error := range errors {
		errorBody[error.Field] = error.Message
	}

	b, _ := json.Marshal(errorBody)
	w.Write(b)
}

func apiAuthenticateUser(r *http.Request) *User {
	if len(r.Header["X-Auth-Token"]) != 1 {
		return nil
	}

	token := r.Header["X-Auth-Token"][0]
	return findUserByAuthToken(token)
}
