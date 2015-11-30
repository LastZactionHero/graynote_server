package main

import "database/sql"

//_, err = db.Exec("CREATE TABLE IF NOT EXISTS notes (id integer AUTO_INCREMENT NOT NULL PRIMARY KEY, user_id integer, title varchar(255), body text)")

// Note stores user note
type Note struct {
	ID     int
	UserID int
	Title  string
	Body   string
}

func createNote(user *User, title string, body string) *Note {
	stmt, err := db.Prepare("INSERT notes SET user_id=?, title=?, body=?")
	checkErr(err, "prepare create note")

	res, err := stmt.Exec(user.ID, title, body)
	checkErr(err, "create note")

	noteID, _ := res.LastInsertId()
	return findNoteByID(noteID)
}

func findNoteByID(noteID int64) *Note {
	var note *Note

	rows, _ := db.Query("SELECT * FROM notes WHERE id=?", noteID)
	if rows.Next() {
		note = noteFromDbRows(rows)
	}
	return note
}

func findNotesByUser(user *User) []*Note {
	rows, _ := db.Query("SELECT * FROM notes WHERE user_id=?", user.ID)
	var notes []*Note
	for rows.Next() {
		notes = append(notes, noteFromDbRows(rows))
	}
	return notes
}

func noteFromDbRows(rows *sql.Rows) *Note {
	note := new(Note)
	rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Body)
	return note
}

// Update a note in the database
func (n *Note) Update(title string, body string) {
	n.Title = title
	n.Body = body

	stmt, err := db.Prepare("UPDATE notes SET title=?, body=? WHERE id=?")
	checkErr(err, "prepare update note")

	_, err = stmt.Exec(title, body, n.ID)
	checkErr(err, "exec update note")
}

// Destroy deletes a Note from the database
func (n Note) Destroy() {
	stmt, err := db.Prepare("DELETE FROM notes WHERE id=?")
	checkErr(err, "prepare delete note")
	_, err = stmt.Exec(n.ID)
	checkErr(err, "exec delete note")
}
