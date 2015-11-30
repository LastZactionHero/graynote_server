package main

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
)

// User account type
type User struct {
	ID           int
	Email        string
	PasswordHash string
	AuthToken    string
}

func passwordToHash(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func createUser(userParams *UserRegisterForm) *User {
	passwordHash := passwordToHash(userParams.Password)
	stmt, _ := db.Prepare("INSERT users SET email=?, password_hash=?, auth_token=?")
	res, err := stmt.Exec(userParams.Email, passwordHash, randomToken())
	checkErr(err, "createUser exec")
	userID, _ := res.LastInsertId()
	user := findUserByID(userID)
	return user
}

func findUserByID(userID int64) *User {
	rows, _ := db.Query("SELECT * FROM users WHERE id=?", userID)

	var user *User

	if rows.Next() {
		user = userFromDbRow(rows)
	}
	return user
}

func findUserByEmail(email string) *User {
	rows, _ := db.Query("SELECT * FROM users WHERE email=?", email)

	var user *User

	if rows.Next() {
		user = userFromDbRow(rows)
	}
	return user
}

func findUserByAuthToken(token string) *User {
	rows, _ := db.Query("SELECT * FROM users WHERE auth_token=?", token)

	var user *User

	if rows.Next() {
		user = userFromDbRow(rows)
	}
	return user
}

func userFromDbRow(rows *sql.Rows) *User {
	user := new(User)
	rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.AuthToken)
	return user
}

func (u User) validPasswordForUser(password string) bool {
	return passwordToHash(password) == u.PasswordHash
}

func randomToken() string {
	b := make([]byte, 64)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
