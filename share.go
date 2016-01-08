package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// Share authenticates sharing of a note
type Share struct {
	ID          int
	NoteID      int
	AuthKey     string
	Permissions string
}

// ValidateSharePermission returns if permission string is valid
func ValidateSharePermission(permissions string) bool {
	return permissions == "readwrite" || permissions == "read"
}

func createShare(note *Note, permissions string) *Share {
	if !ValidateSharePermission(permissions) {
		return nil
	}

	stmt, err := db.Prepare("INSERT shares SET note_id=?, auth_key=?, permissions=?")
	if err != nil {
		checkErr(err, "prepare create share")
	} else {
		defer stmt.Close()
	}
	res, err := stmt.Exec(note.ID, randomShareKey(), permissions)
	checkErr(err, "create share")

	shareID, _ := res.LastInsertId()
	return findShareByID(shareID)
}

func findShareByID(id int64) *Share {
	var share *Share
	rows, err := db.Query("SELECT * FROM shares WHERE id=?", id)
	if err != nil {
		checkErr(err, "find share by id")
	} else {
		defer rows.Close()
	}

	if rows.Next() {
		share = shareFromDbRows(rows)
	}

	return share
}

func findShareByAuthKey(authKey string) *Share {
	var share *Share
	rows, err := db.Query("SELECT * FROM shares WHERE auth_key=?", authKey)
	if err != nil {
		checkErr(err, "find share by auth key")
	} else {
		defer rows.Close()
	}

	if rows.Next() {
		share = shareFromDbRows(rows)
	}

	return share
}

// Destoy a share from database
func (s Share) Destroy() {
	stmt, err := db.Prepare("DELETE FROM shares WHERE id=?")
	if err != nil {
		checkErr(err, "prepare destroy")
	} else {
		defer stmt.Close()
	}

	_, err = stmt.Exec(s.ID)
	checkErr(err, "delete share")
}

func shareFromDbRows(rows *sql.Rows) *Share {
	share := new(Share)
	rows.Scan(&share.ID, &share.AuthKey, &share.NoteID, &share.Permissions)
	return share
}

// TODO: IMPROVE RANDOM KEY
func randomShareKey() string {
	hasher := md5.New()
	rand.Seed(int64(time.Now().Nanosecond()))
	rand := rand.Int31()
	hasher.Write([]byte(fmt.Sprintf("%d", rand)))
	return hex.EncodeToString(hasher.Sum(nil))
}
