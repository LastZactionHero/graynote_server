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
	var err error

	dbUser := os.Getenv("GRAYNOTE_DB_USER")
	dbPass := os.Getenv("GRAYNOTE_DB_PASS")
	dbName := os.Getenv("GRAYNOTE_DB_NAME")
	dbConnect := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPass, dbName)

	db, err = sql.Open("mysql", dbConnect)
	db.SetMaxIdleConns(10000)
	db.SetMaxOpenConns(10000)
	checkErr(err, "sql.Open failed")
	defer db.Close()

	err = db.Ping()
	checkErr(err, "db ping failed")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS notes (id integer AUTO_INCREMENT NOT NULL PRIMARY KEY, user_id integer, title varchar(255), body text)")
	checkErr(err, "create table Notes failed")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id integer AUTO_INCREMENT NOT NULL PRIMARY KEY, email varchar(255), password_hash varchar(255), auth_token varchar(64))")
	checkErr(err, "create table Users failed")

	r := mux.NewRouter()
	r.HandleFunc("/notes", optionsHandler).Methods("OPTIONS")
	r.HandleFunc("/notes/{id:[0-9]+}", optionsHandler).Methods("OPTIONS")

	r.HandleFunc("/users/register", userRegisterHandler).Methods("POST")
	r.HandleFunc("/users/login", userLoginHandler).Methods("POST")

	r.HandleFunc("/notes", noteIndexHandler).Methods("GET")
	r.HandleFunc("/notes", noteCreateHandler).Methods("POST")
	r.HandleFunc("/notes/{id:[0-9]+}", noteShowHandler).Methods("GET")
	r.HandleFunc("/notes/{id:[0-9]+}", noteUpdateHandler).Methods("PUT")
	r.HandleFunc("/notes/{id:[0-9]+}", noteDeleteHandler).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8181", nil)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
