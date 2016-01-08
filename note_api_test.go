package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestNoteCreateHandlerSuccess(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	title := "My Note!"
	body := "Some exciting things are documented here."

	// Create Notes
	postBody := strings.NewReader(fmt.Sprintf("title=%s&body=%s", title, body))
	r, _ := http.NewRequest("POST", "/notes", postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 201 {
		t.Errorf("Expected 201, got %q", w.Code)
	}

	expectedBody := "{\"id\":1,\"title\":\"My Note!\",\"body\":\"Some exciting things are documented here.\",\"shares\":null}"
	if b := w.Body.String(); b != expectedBody {
		t.Errorf("Expected %q, but got %q", expectedBody, b)
	}
}

func TestNoteCreateFailMissingParams(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create Notes
	postBody := strings.NewReader(fmt.Sprintf(""))
	r, _ := http.NewRequest("POST", "/notes", postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, got %q", w.Code)
	}

	expectedBody := "{\"body\":\"is required\",\"title\":\"is required\"}"
	if b := w.Body.String(); b != expectedBody {
		t.Errorf("Expected %q, but got %q", expectedBody, b)
	}
}

func TestNoteDeleteHandlerSuccess(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create a note
	note := createNote(user, "My Note", "Note Body!")

	// Delete Note
	path := fmt.Sprintf("/notes/%d", note.ID)
	r, _ := http.NewRequest("DELETE", path, nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %q", w.Code)
	}

	if findNoteByID(int64(note.ID)) != nil {
		t.Errorf("Expected note to be deleted")
	}
}

func TestNoteDeleteHandlerFailInvalidUser(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create a note
	otherUser := factoryCreateUser("someone@else.com")
	note := createNote(otherUser, "My Note", "Note Body!")

	// Delete Note
	path := fmt.Sprintf("/notes/%d", note.ID)
	r, _ := http.NewRequest("DELETE", path, nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %q", w.Code)
	}

	if findNoteByID(int64(note.ID)) == nil {
		t.Errorf("Expected note not to be deleted")
	}
}

func TestNoteDeleteHandlerFailNotFound(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Delete Note
	path := "/notes/1"
	r, _ := http.NewRequest("DELETE", path, nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %q", w.Code)
	}
}

func TestNoteIndexHandler(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a user
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create notes
	createNote(user, "My Note", "Note Body!")
	createNote(user, "Second Note", "Second Note Body!")

	// Get Notes index
	r, _ := http.NewRequest("GET", "/notes", nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected code 200, got %q", w.Code)
	}
	expectedBody := "[{\"id\":1,\"title\":\"My Note\",\"body\":\"Note Body!\",\"shares\":null},{\"id\":2,\"title\":\"Second Note\",\"body\":\"Second Note Body!\",\"shares\":null}]"
	if b := w.Body.String(); b != expectedBody {
		t.Errorf("Expected %q, got %q", expectedBody, b)
	}
}

func TestNoteShowHandlerSuccess(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create a note
	note := createNote(user, "My Note", "Note Body!")

	// Get the note
	path := fmt.Sprintf("/notes/%d", note.ID)
	r, _ := http.NewRequest("GET", path, nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %q", w.Code)
	}

	expectedBody := "{\"id\":1,\"title\":\"My Note\",\"body\":\"Note Body!\",\"shares\":null}"
	if b := w.Body.String(); b != expectedBody {
		t.Errorf("Expected %q, got %q", expectedBody, b)
	}
}

func TestNoteShowHandlerSuccessWithShares(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create a note
	note := createNote(user, "My Note", "Note Body!")

	createShare(note, "readwrite")
	createShare(note, "read")

	// Get the note
	path := fmt.Sprintf("/notes/%d", note.ID)
	r, _ := http.NewRequest("GET", path, nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %q", w.Code)
	}

	b := w.Body.String()
	expected := fmt.Sprintf(
		"{\"id\":1,\"title\":\"My Note\",\"body\":\"Note Body!\",\"shares\":[{\"auth_key\":\"[a-f0-9]+\",\"note_id\":1,\"permissions\":\"readwrite\"},{\"auth_key\":\"[a-f0-9]+\",\"note_id\":1,\"permissions\":\"read\"}]}")
	if match, _ := regexp.MatchString(expected, b); !match {
		t.Errorf("Expected %q to match %q", b, expected)
	}
}

func TestNoteShowHandlerFailDoesNotExist(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Get the now
	r, _ := http.NewRequest("GET", "/notes/1", nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %q", w.Code)
	}
}

func TestNoteShowHandlerFailDifferentOwner(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create a note for a difference user
	otherUser := factoryCreateUser("someone@else.com")
	note := createNote(otherUser, "My Note", "Note Body!")

	// Get the now
	path := fmt.Sprintf("/notes/%d", note.ID)
	r, _ := http.NewRequest("GET", path, nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %q", w.Code)
	}
}

func TestNoteUpdateHandlerSuccess(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create a note
	note := createNote(user, "My Note", "Note Body!")

	title := "Updated Title"
	body := "Updated Body"

	// Update Notes
	postBody := strings.NewReader(fmt.Sprintf("title=%s&body=%s", title, body))
	path := fmt.Sprintf("/notes/%d", note.ID)
	r, _ := http.NewRequest("PUT", path, postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %q", w.Code)
	}
	expectedBody := "{\"id\":1,\"title\":\"Updated Title\",\"body\":\"Updated Body\",\"shares\":null}"
	if b := w.Body.String(); b != expectedBody {
		t.Errorf("Expected %q, got %q", expectedBody, b)
	}

	note = findNoteByID(int64(note.ID))
	if note.Title != title || note.Body != body {
		t.Errorf("Expected note to equal updated values")
	}
}

func TestNoteUpdateHandlerFailMissingParameters(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create a note
	note := createNote(user, "My Note", "Note Body!")

	// Update Notes
	postBody := strings.NewReader("")
	path := fmt.Sprintf("/notes/%d", note.ID)
	r, _ := http.NewRequest("PUT", path, postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, got %q", w.Code)
	}
	expectedBody := "{\"body\":\"is required\",\"title\":\"is required\"}"
	if b := w.Body.String(); b != expectedBody {
		t.Errorf("Expected %q, got %q", expectedBody, b)
	}
}

func TestNoteUpdateHandlerFailNotFound(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Update Notes
	postBody := strings.NewReader("")
	path := "/notes/1"
	r, _ := http.NewRequest("PUT", path, postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %q", w.Code)
	}
}

func TestNoteUpdateHandlerFailInvalidOwner(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	// Create a User
	userEmail := "user@site.com"
	user := factoryCreateUser(userEmail)

	// Create a note
	otherUser := factoryCreateUser("someone@else.com")
	note := createNote(otherUser, "My Note", "Note Body!")

	title := "Updated Title"
	body := "Updated Body"

	// Update Notes
	postBody := strings.NewReader(fmt.Sprintf("title=%s&body=%s", title, body))
	path := fmt.Sprintf("/notes/%d", note.ID)
	r, _ := http.NewRequest("PUT", path, postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %q", w.Code)
	}
}
