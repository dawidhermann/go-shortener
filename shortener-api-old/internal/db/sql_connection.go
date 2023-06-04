package db

import (
	"database/sql"
	"fmt"
	"github.com/dawidhermann/shortener-api/config"
	"log"
)

type SqlConnection struct {
	Db *sql.DB
}

func Connect(dbConfig config.DbConfig) SqlConnection {
	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", dbConfig.DbUser, dbConfig.DbPassword, dbConfig.DbAddr, dbConfig.DbName)
	fmt.Println(connString)
	dbPtr, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
	}
	return SqlConnection{Db: dbPtr}
}

func (sqlConn SqlConnection) Close() {
	sqlConn.Close()
}
