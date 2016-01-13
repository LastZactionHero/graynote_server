package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var sessionStore = sessions.NewCookieStore([]byte(os.Getenv("GRAYNOTE_SESSION_KEY")))

func main() {
	fmt.Println("Graynote Server")
	dbSetup(
		os.Getenv("GRAYNOTE_DB_USER"),
		os.Getenv("GRAYNOTE_DB_PASS"),
		os.Getenv("GRAYNOTE_DB_NAME"),
		false)

	defer db.Close()

	r := router()
	http.Handle("/", r)
	http.ListenAndServe(":8181", nil)
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/notes", optionsHandler).Methods("OPTIONS")
	r.HandleFunc("/notes/{id:[a-z0-9]+}", optionsHandler).Methods("OPTIONS")
	r.HandleFunc("/shares", optionsHandler).Methods("OPTIONS")
	r.HandleFunc("/shares/{id:[a-z0-9]+}", optionsHandler).Methods("OPTIONS")

	r.HandleFunc("/users/register", userRegisterHandler).Methods("POST")
	r.HandleFunc("/users/login", userLoginHandler).Methods("POST")

	r.HandleFunc("/notes", noteIndexHandler).Methods("GET")
	r.HandleFunc("/notes", noteCreateHandler).Methods("POST")
	r.HandleFunc("/notes/{id:[a-z0-9]+}", noteShowHandler).Methods("GET")
	r.HandleFunc("/notes/{id:[a-z0-9]+}", noteUpdateHandler).Methods("PUT")
	r.HandleFunc("/notes/{id:[0-9]+}", noteDeleteHandler).Methods("DELETE")

	r.HandleFunc("/shares", shareCreateHandler).Methods("POST")
	r.HandleFunc("/shares/{id:[a-z0-9]+}", shareDeleteHandler).Methods("DELETE")

	return r
}

func dbSetup(dbUser string, dbPass string, dbName string, wipe bool) {
	var err error

	dbConnect := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPass, dbName)

	db, err = sql.Open("mysql", dbConnect)
	db.SetMaxIdleConns(10000)
	db.SetMaxOpenConns(10000)
	checkErr(err, "sql.Open failed")

	err = db.Ping()
	checkErr(err, "db ping failed")

	if wipe {
		_, err = db.Exec("DROP TABLE IF EXISTS users")
		checkErr(err, "drop table users")
		_, err = db.Exec("DROP TABLE IF EXISTS notes")
		checkErr(err, "drop table notes")
		_, err = db.Exec("DROP TABLE IF EXISTS shares")
		checkErr(err, "drop table shares")
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS notes (id integer AUTO_INCREMENT NOT NULL PRIMARY KEY, user_id integer, title varchar(255), body text)")
	checkErr(err, "create table Notes failed")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id integer AUTO_INCREMENT NOT NULL PRIMARY KEY, email varchar(255), password_hash varchar(255), auth_token varchar(64))")
	checkErr(err, "create table Users failed")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS shares (id integer AUTO_INCREMENT NOT NULL PRIMARY KEY, auth_key varchar(255), note_id integer, permissions varchar(255))")
	checkErr(err, "create table Shares failed")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
