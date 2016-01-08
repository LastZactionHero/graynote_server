package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type shareRequestParameters struct {
	NoteID      int    `schema:"note_id"`
	Permissions string `schema:"permissions"`
}

type shareSuccessResponse struct {
	AuthKey     string `json:"auth_key"`
	NoteID      int    `json:"note_id"`
	Permissions string `json:"permissions"`
}

func shareCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	apiApplyCorsHeaders(w, r)

	// Authenticate
	user := apiAuthenticateUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	checkErr(err, "parsing form")

	shareParameters := new(shareRequestParameters)
	decoder := schema.NewDecoder()
	err = decoder.Decode(shareParameters, r.PostForm)
	checkErr(err, "decoding share create")
	note := findNoteByID(int64(shareParameters.NoteID))

	// Validate Note Exists and is owned by User
	if note == nil || note.UserID != user.ID {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var errors []APIError

	// Valdiate Permissions
	if len(shareParameters.Permissions) <= 0 {
		errors = append(errors, APIError{Field: "permissions", Message: "is required"})
	} else if !ValidateSharePermission(shareParameters.Permissions) {
		errors = append(errors, APIError{Field: "permissions", Message: "is invalid"})
	}

	if len(errors) > 0 {
		apiErrorHandler(w, r, http.StatusBadRequest, errors)
		return
	}

	// Create Share
	share := createShare(note, shareParameters.Permissions)

	// Success message
	w.WriteHeader(http.StatusCreated)
	w.Write(shareJSON(share))
}

func shareDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	apiApplyCorsHeaders(w, r)

	// Authenticate
	user := apiAuthenticateUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Fetch Share by AuthKey
	shareAuthKey := mux.Vars(r)["id"]
	share := findShareByAuthKey(shareAuthKey)

	// Validate Share Exists
	if share == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find Note for Share
	note := findNoteByID(int64(share.NoteID))

	// Validate Note belongs to User
	if note == nil || note.UserID != user.ID {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete the Share
	share.Destroy()

	w.Write([]byte("{}"))
}

func shareJSON(share *Share) []byte {
	response := shareSuccessResponse{AuthKey: share.AuthKey, NoteID: share.NoteID, Permissions: share.Permissions}
	responseJSON, _ := json.Marshal(response)
	return responseJSON
}
