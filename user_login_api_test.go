package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserRegisterHandlerSuccess(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Test successful create
	postBody := "email=user@site.com&password=mypassword.com"
	postBodyReader := strings.NewReader(postBody)

	r, _ := http.NewRequest("POST", "/users/register", postBodyReader)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 201 {
		t.Errorf("Expected code 201, got %q", w.Code)
	}
	if b := w.Body.String(); !strings.Contains(b, "token") {
		t.Errorf("Expected token, got %q", b)
	}

	user := findUserByEmail("user@site.com")
	if user == nil {
		t.Errorf("Expected user, got nil")
	}
}

func TestUserRegisterHandlerFailMissingParams(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	postBody := ""
	postBodyReader := strings.NewReader(postBody)

	r, _ := http.NewRequest("POST", "/users/register", postBodyReader)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected code 400, got %q", w.Code)
	}
	expectedError := "{\"email\":\"is required\",\"password\":\"is required\"}"
	if b := w.Body.String(); b != expectedError {
		t.Errorf("Expected %q, got %q", expectedError, b)
	}
}

func TestUserRegisterHandlerFailAlreadyExists(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	form := UserRegisterForm{Email: "user@site.com", Password: "password"}
	createUser(&form)

	postBody := "email=user@site.com&password=mypassword.com"
	postBodyReader := strings.NewReader(postBody)

	// Create the first user
	r, _ := http.NewRequest("POST", "/users/register", postBodyReader)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected code 400, got %q", w.Code)
	}

	expectedError := "{\"email\":\"already exists\"}"
	if b := w.Body.String(); b != expectedError {
		t.Errorf("Expected %q, got %q", expectedError, b)
	}
}

func TestUserLoginHandlerSuccess(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create the user
	form := UserRegisterForm{Email: "user@site.com", Password: "thepassword"}
	createUser(&form)

	postBody := "email=user@site.com&password=thepassword"
	postBodyReader := strings.NewReader(postBody)

	r, _ := http.NewRequest("POST", "/users/login", postBodyReader)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 201 {
		t.Errorf("Expected code 201, got %q", w.Code)
	}
	if b := w.Body.String(); !strings.Contains(b, "token") {
		t.Errorf("Expected token, got %q", b)
	}
}

func TestUserLoginHandlerFailBadPassword(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create the user
	form := UserRegisterForm{Email: "user@site.com", Password: "thepassword"}
	createUser(&form)

	postBody := "email=user@site.com&password=wrongpassword"
	postBodyReader := strings.NewReader(postBody)

	r, _ := http.NewRequest("POST", "/users/login", postBodyReader)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected code 403, got %q", w.Code)
	}
	expectedError := "{\"password\":\"is invalid\"}"
	if b := w.Body.String(); b != expectedError {
		t.Errorf("Expected %q, got %q", expectedError, b)
	}
}

func TestUserLoginHandlerFailNoEmail(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create the user
	form := UserRegisterForm{Email: "user@site.com", Password: "thepassword"}
	createUser(&form)

	postBody := "email=missinguser@site.com&password=wrongpassword"
	postBodyReader := strings.NewReader(postBody)

	r, _ := http.NewRequest("POST", "/users/login", postBodyReader)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected code 400, got %q", w.Code)
	}
	expectedError := "{\"email\":\"not found\"}"
	if b := w.Body.String(); b != expectedError {
		t.Errorf("Expected %q, got %q", expectedError, b)
	}
}
