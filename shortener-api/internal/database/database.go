package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/dawidhermann/shortener-api/appbase/config"
	"github.com/jmoiron/sqlx"
)

var (
	ErrDbNotFound       = errors.New("db record not found")
	ErrDoNoRowsAffected = errors.New("no rows affected")
)

func Connect(cfg config.DbConfig) (*sqlx.DB, error) {
	q := make(url.Values)
	q.Set("sslmode", "disable")
	connectionUrl := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:     cfg.DbName,
		RawQuery: q.Encode(),
	}
	fmt.Println(connectionUrl.String())
	dbPtr, err := sqlx.Open("postgres", connectionUrl.String())
	if err != nil {
		return nil, err
	}
	if err = dbPtr.Ping(); err != nil {
		return nil, err
	}
	return dbPtr, nil
}

func NamedQueryStruct(ctx context.Context, db sqlx.ExtContext, query string, queryData any, queryDest any) error {
	rows, err := sqlx.NamedQueryContext(ctx, db, query, queryData)
	if err != nil {
		return fmt.Errorf("failed to execute query, error: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return ErrDbNotFound
	}
	if err = rows.StructScan(queryDest); err != nil {
		return fmt.Errorf("failed to scan data, error: %w", err)
	}
	return nil
}

func NamedExecContext(ctx context.Context, db sqlx.ExtContext, query string, queryData any) error {
	result, err := sqlx.NamedExecContext(ctx, db, query, queryData)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrDoNoRowsAffected
	}
	return nil
}
