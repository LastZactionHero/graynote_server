package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type noteRequestParameters struct {
	Title string `schema:"title"`
	Body  string `schema:"body"`
}

type noteSuccessResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func noteIndexHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate
	user := apiAuthenticateUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	notes := findNotesByUser(user)
	w.Write(notesJSON(notes))
}

func noteCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate
	user := apiAuthenticateUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	checkErr(err, "parsing form")

	noteParameters := new(noteRequestParameters)
	decoder := schema.NewDecoder()
	err = decoder.Decode(noteParameters, r.PostForm)
	checkErr(err, "decoding note create")

	var errors []APIError

	// Validate Title
	if len(noteParameters.Title) == 0 {
		errors = append(errors, APIError{Field: "title", Message: "is required"})
	}

	// Validate Body
	if len(noteParameters.Body) == 0 {
		errors = append(errors, APIError{Field: "body", Message: "is required"})
	}

	if len(errors) > 0 {
		apiErrorHandler(w, r, http.StatusBadRequest, errors)
		return
	}

	// Create Note
	note := createNote(user, noteParameters.Title, noteParameters.Body)

	// Success message
	w.WriteHeader(http.StatusCreated)
	w.Write(noteJSON(note))
}

func noteShowHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate
	user := apiAuthenticateUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Find the note
	noteIDStr := mux.Vars(r)["id"]
	noteID, _ := strconv.ParseInt(noteIDStr, 10, 64)
	note := findNoteByID(noteID)

	// Note not found or invalid owner
	if note == nil || note.UserID != user.ID {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write(noteJSON(note))
}

func noteUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate
	user := apiAuthenticateUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Find the note
	noteIDStr := mux.Vars(r)["id"]
	noteID, _ := strconv.ParseInt(noteIDStr, 10, 64)
	note := findNoteByID(noteID)

	// Note not found or invalid owner
	if note == nil || note.UserID != user.ID {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	checkErr(err, "parsing form")

	noteParameters := new(noteRequestParameters)
	decoder := schema.NewDecoder()
	err = decoder.Decode(noteParameters, r.PostForm)
	checkErr(err, "decoding note create")

	var errors []APIError

	// Validate Title
	if len(noteParameters.Title) == 0 {
		errors = append(errors, APIError{Field: "title", Message: "is required"})
	}

	// Validate Body
	if len(noteParameters.Body) == 0 {
		errors = append(errors, APIError{Field: "body", Message: "is required"})
	}

	if len(errors) > 0 {
		apiErrorHandler(w, r, http.StatusBadRequest, errors)
		return
	}

	note.Update(noteParameters.Title, noteParameters.Body)

	w.Write(noteJSON(note))
}

func noteDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate
	user := apiAuthenticateUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Find the note
	noteIDStr := mux.Vars(r)["id"]
	noteID, _ := strconv.ParseInt(noteIDStr, 10, 64)
	note := findNoteByID(noteID)

	// Note not found or invalid owner
	if note == nil || note.UserID != user.ID {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	note.Destroy()
}

func noteJSON(note *Note) []byte {
	response := noteSuccessResponse{ID: note.ID, Title: note.Title, Body: note.Body}
	responseJSON, _ := json.Marshal(response)
	return responseJSON
}

func notesJSON(notes []*Note) []byte {
	var response []noteSuccessResponse
	for _, note := range notes {
		response = append(response,
			noteSuccessResponse{ID: note.ID, Title: note.Title, Body: note.Body})
	}
	responseJSON, _ := json.Marshal(response)
	return responseJSON
}
