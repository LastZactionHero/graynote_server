package main

import (
	"database/sql"
	"os"
)

func testDbSetup() *sql.DB {
	dbSetup(
		os.Getenv("GRAYNOTE_DB_USER"),
		os.Getenv("GRAYNOTE_DB_PASS"),
		os.Getenv("GRAYNOTE_DB_TEST_NAME"),
		true)
	return db
}

func factoryCreateUser(email string) *User {
	form := UserRegisterForm{Email: email, Password: "password"}
	createUser(&form)
	return findUserByEmail(email)
}
