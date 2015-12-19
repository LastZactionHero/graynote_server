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

// TODO: IMPROVE PASSWORD HASH
func passwordToHash(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func createUser(userParams *UserRegisterForm) *User {
	passwordHash := passwordToHash(userParams.Password)
	stmt, _ := db.Prepare("INSERT users SET email=?, password_hash=?, auth_token=?")
	fmt.Println(randomToken())
	res, err := stmt.Exec(userParams.Email, passwordHash, randomToken())
	checkErr(err, "createUser exec")
	userID, _ := res.LastInsertId()
	user := findUserByID(userID)
	return user
}

func findUserByID(userID int64) *User {
	rows, err := db.Query("SELECT * FROM users WHERE id=?", userID)
	if err != nil {
		checkErr(err, "findUserByID")
	} else {
		defer rows.Close()
	}

	var user *User

	if rows.Next() {
		user = userFromDbRow(rows)
	}
	return user
}

func findUserByEmail(email string) *User {
	rows, err := db.Query("SELECT * FROM users WHERE email=?", email)
	if err != nil {
		checkErr(err, "findUserByEmail")
	} else {
		defer rows.Close()
	}

	var user *User

	if rows.Next() {
		user = userFromDbRow(rows)
	}
	return user
}

func findUserByAuthToken(token string) *User {
	rows, err := db.Query("SELECT * FROM users WHERE auth_token=?", token)
	if err != nil {
		checkErr(err, "findUserByAuthToken")
	} else {
		defer rows.Close()
	}

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
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
