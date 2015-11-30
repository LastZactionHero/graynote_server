package main

// UserRegisterForm type
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

// UserRegisterForm parameters for login and registration
type UserRegisterForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type userLoginSuccessResponse struct {
	Token string `json:"token"`
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle Error
		checkErr(err, "error parsing form")
	}

	userParams := new(UserRegisterForm)
	decoder := schema.NewDecoder()
	err = decoder.Decode(userParams, r.PostForm)

	// Error message if missing email and/or password
	var paramErrors []APIError
	if len(userParams.Email) == 0 {
		error := APIError{Field: "email", Message: "is required"}
		paramErrors = append(paramErrors, error)
	}

	if len(userParams.Password) == 0 {
		error := APIError{Field: "password", Message: "is required"}
		paramErrors = append(paramErrors, error)
	}

	if len(paramErrors) > 0 {
		apiErrorHandler(w, r, http.StatusBadRequest, paramErrors)
		return
	}

	// Error message if user already exists
	user := findUserByEmail(userParams.Email)
	if user != nil {
		error := APIError{Field: "email", Message: "already exists"}
		apiErrorHandler(w, r, http.StatusBadRequest, []APIError{error})
		return
	}

	// Create user
	user = createUser(userParams)
	fmt.Println(user.Email)

	// Success message
	w.WriteHeader(http.StatusCreated)

	// Authenticate
	w.Write(successfulLoginJSON(user))
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("userLoginHandler")
	err := r.ParseForm()
	if err != nil {
		// Handle Error
		checkErr(err, "error parsing form")
	}

	userParams := new(UserRegisterForm)
	decoder := schema.NewDecoder()
	err = decoder.Decode(userParams, r.PostForm)

	var paramErrors []APIError

	// Error if email is not present
	if len(userParams.Email) == 0 {
		paramErrors = append(paramErrors, APIError{Field: "email", Message: "is required"})
	}

	// Error if password is not present
	if len(userParams.Password) == 0 {
		paramErrors = append(paramErrors, APIError{Field: "password", Message: "is required"})
	}

	if len(paramErrors) > 0 {
		apiErrorHandler(w, r, http.StatusBadRequest, paramErrors)
		return
	}

	// Check for user
	user := findUserByEmail(userParams.Email)

	// Error if user is not found
	if user == nil {
		error := APIError{Field: "email", Message: "not found"}
		apiErrorHandler(w, r, http.StatusForbidden, []APIError{error})
		return
	}

	// Validate password for user
	// Error if password is invalid
	if !user.validPasswordForUser(userParams.Password) {
		error := APIError{Field: "password", Message: "is invalid"}
		apiErrorHandler(w, r, http.StatusForbidden, []APIError{error})
		return
	}

	// Success message
	w.WriteHeader(http.StatusCreated)

	// Authenticate
	w.Write(successfulLoginJSON(user))
}

func successfulLoginJSON(user *User) []byte {
	response := userLoginSuccessResponse{Token: user.AuthToken}
	responseJSON, _ := json.Marshal(response)
	return responseJSON
}
