package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestShareCreateHandlerSuccess(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")

	permissions := "readwrite"
	postBody := strings.NewReader(fmt.Sprintf("note_id=%d&permissions=%s", note.ID, permissions))
	r, _ := http.NewRequest("POST", "/shares", postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 201 {
		t.Errorf("Expected 201 response, got %q", w.Code)
	}

	b := w.Body.String()
	expected := fmt.Sprintf(
		"{\"auth_key\":\"[a-f0-9]+\",\"note_id\":%d,\"permissions\":\"%s\"}",
		note.ID,
		permissions)
	if match, _ := regexp.MatchString(expected, b); !match {
		t.Errorf("Expected %q to match %q", b, expected)
	}
}

func TestShareCreateHandlerFailInvalidUser(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")

	otherUser := factoryCreateUser("other@guy.com")
	note := createNote(otherUser, "title", "body")

	permissions := "readwrite"
	postBody := strings.NewReader(fmt.Sprintf("note_id=%d&permissions=%s", note.ID, permissions))
	r, _ := http.NewRequest("POST", "/shares", postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404 response, got %q", w.Code)
	}
}

func TestShareCreateHandlerFailInvalidPermissions(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")

	permissions := "garbage"
	postBody := strings.NewReader(fmt.Sprintf("note_id=%d&permissions=%s", note.ID, permissions))
	r, _ := http.NewRequest("POST", "/shares", postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400 response, got %q", w.Code)
	}

	expectedBody := "{\"permissions\":\"is invalid\"}"
	if b := w.Body.String(); b != expectedBody {
		t.Errorf("Expected %q, got %q", expectedBody, b)
	}
}

func TestShareCreateHandlerFailMissingPermissions(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")

	postBody := strings.NewReader(fmt.Sprintf("note_id=%d", note.ID))
	r, _ := http.NewRequest("POST", "/shares", postBody)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400 response, got %q", w.Code)
	}

	expectedBody := "{\"permissions\":\"is required\"}"
	if b := w.Body.String(); b != expectedBody {
		t.Errorf("Expected %q, got %q", expectedBody, b)
	}
}

func TestShareDeleteHandlerSuccess(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")
	share := createShare(note, "readwrite")

	path := fmt.Sprintf("/shares/%s", share.AuthKey)
	r, _ := http.NewRequest("DELETE", path, nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200 response, got %q", w.Code)
	}
}

func TestShareDeleteHandlerInvalidUser(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")

	otherUser := factoryCreateUser("other@guy.com")
	note := createNote(otherUser, "title", "body")
	share := createShare(note, "readwrite")

	path := fmt.Sprintf("/shares/%s", share.AuthKey)
	r, _ := http.NewRequest("DELETE", path, nil)
	r.Header.Add("X-Auth-Token", user.AuthToken)
	w := httptest.NewRecorder()

	router().ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Expected 404 response, got %q", w.Code)
	}
}
