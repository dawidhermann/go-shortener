package store

import (
	"context"
	"errors"
	"time"

	"github.com/dawidhermann/shortener-api/internal/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	uniqueViolation = "23505"
	undefinedTable  = "42P01"
)

var (
	ErrUndefinedTable       = errors.New("undefined table")
	ErrUniqueViolation      = errors.New("unique violation")
	ErrUniqueEmailViolation = errors.New("unique email violation")
	ErrUserNotFound         = errors.New("user not found")
)

type Store struct {
	db sqlx.ExtContext
}

type CreateUrlResult struct {
	UrlId     uuid.UUID `db:"url_id"`
	CreatedAt time.Time `db:"created_at"`
}

func NewUrlStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (store Store) Create(ctx context.Context, url DbUrl) (CreateUrlResult, error) {
	const query = `
			INSERT INTO urls (url_key, user_id)
			VALUES (:url_key, :url_key)
			RETURNING url_id, created_at
		`
	var createUrlResult CreateUrlResult
	if err := database.NamedQueryStruct(ctx, store.db, query, url, &createUrlResult); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case undefinedTable:
				return CreateUrlResult{}, ErrUndefinedTable
			case uniqueViolation:
				return CreateUrlResult{}, ErrUniqueViolation
			}
		}
		return createUrlResult, err
	}
	return createUrlResult, nil
}

func (store Store) GetById(ctx context.Context, id uuid.UUID) (DbUrl, error) {
	queryData := struct {
		UrlId uuid.UUID `db:"url_id"`
	}{
		UrlId: id,
	}
	const query = `
		SELECT
		"url_id",
		"url_key",
		"user_id",
		"created_at"
		FROM urls
		WHERE url_id = :url_id
	`
	var dbUrl DbUrl
	if err := database.NamedQueryStruct(ctx, store.db, query, queryData, &dbUrl); err != nil {
		if errors.Is(err, database.ErrDbNotFound) {
			return DbUrl{}, ErrUserNotFound
		}
		return DbUrl{}, err
	}
	return dbUrl, nil
}

func (store Store) Delete(ctx context.Context, id uuid.UUID) error {
	queryData := struct {
		UrlId uuid.UUID `db:"url_id"`
	}{
		UrlId: id,
	}
	const query = `
		DELETE from urls
		WHERE url_id = :url_id
	`
	if err := database.NamedExecContext(ctx, store.db, query, queryData); err != nil {
		if err == database.ErrDoNoRowsAffected {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
