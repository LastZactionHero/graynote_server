package main

import "testing"

func TestCreateNote(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	title := "Note Title"
	body := "Note Body"

	note := createNote(user, title, body)
	if note == nil {
		t.Errorf("Expected note not to be nil")
	}
	if note.Title != title {
		t.Errorf("Expected note title to eq %q", title)
	}
	if note.Body != body {
		t.Errorf("Expected note body to eq %q", body)
	}
	if note.ID != 1 {
		t.Errorf("Expected note id to eq %q, got %q", 1, note.ID)
	}
}

func TestFindNoteByID(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")

	foundNote := findNoteByID(int64(note.ID))
	if foundNote.ID != note.ID || foundNote.Title != note.Title || foundNote.Body != note.Body {
		t.Errorf("Expected note to equal original")
	}
}

func TestFindNotesByUser(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	noteA := createNote(user, "title", "body")
	noteB := createNote(user, "title", "body")
	noteC := createNote(user, "title", "body")

	notes := findNotesByUser(user, "")
	if notes[0].ID != noteA.ID {
		t.Errorf("Expected first note to be note_a")
	}
	if notes[1].ID != noteB.ID {
		t.Errorf("Expected second note to be note_b")
	}
	if notes[2].ID != noteC.ID {
		t.Errorf("Expected third note to be note_c")
	}
}

func TestFindNotesByUserWithQuery(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	createNote(user, "title", "body")
	createNote(user, "title", "body")
	queryNote := createNote(user, "title", "this should match the query")

	notes := findNotesByUser(user, "the query")
	if len(notes) != 1 {
		t.Errorf("Expected to find 1 note, found %d", len(notes))
	}
	if notes[0].ID != queryNote.ID {
		t.Errorf("Expected queryNote to be found")
	}
}

func TestNoteDestroy(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")
	note.Destroy()

	if n := findNoteByID(int64(note.ID)); n != nil {
		t.Errorf("Expected note to not exist, got ID %d", n.ID)
	}
}

func TestNoteUpdate(t *testing.T) {
	db := testDbSetup()
	defer db.Close()

	user := factoryCreateUser("user@site.com")
	note := createNote(user, "title", "body")

	title := "updated title"
	body := "updated body"
	note.Update(title, body)

	updated := findNoteByID(int64(note.ID))
	if updated.Title != title {
		t.Errorf("Expected title to be %q, got %q", title, updated.Title)
	}
	if updated.Body != body {
		t.Errorf("Expected body to be %q, got %q", body, updated.Body)
	}
}
