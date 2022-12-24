package db

import (
	"database/sql"
	"fmt"
	"log"
)

var Db *sql.DB

func Connect(username string, password string, dbHost string, dbName string) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", username, password, dbHost, dbName)
	fmt.Println(connString)
	dbPtr, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
	}
	Db = dbPtr
}
