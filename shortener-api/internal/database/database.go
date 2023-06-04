package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
)

type DbConfig struct {
	User     string
	Password string
	Host     string
	Name     string
	Schema   string
}

func Connect(cfg DbConfig) (*sqlx.DB, error) {
	q := make(url.Values)
	connectionUrl := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
	dbPtr, err := sqlx.Open("postgres", connectionUrl.String())
	if err != nil {
		return nil, err
	}
	return dbPtr, nil
}

//func (sqlConn SqlConnection) Close() {
//	sqlConn.Close()
//}
