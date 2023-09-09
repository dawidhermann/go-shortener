// url db store
package store

import (
	"context"
	"errors"
	"fmt"
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
	ErrUndefinedTable  = errors.New("undefined table")
	ErrUniqueViolation = errors.New("unique violation")
	ErrUrlNotFound     = errors.New("url not found")
)

type TxFunc func(key string) error

type Store struct {
	db sqlx.ExtContext
}

type CreateUrlResult struct {
	UrlId     uuid.UUID `db:"url_id"`
	CreatedAt time.Time `db:"created_at"`
}

// Create new url store
func NewUrlStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// Create new db url
func (store Store) Create(ctx context.Context, url DbUrl) (CreateUrlResult, error) {
	const query = `
			INSERT INTO urls (url_key, user_id)
			VALUES (:url_key, :user_id)
			RETURNING url_id, created_at
		`
	var createUrlResult CreateUrlResult
	// 	tx, err := database.WithTx(store.db)
	// if err != nil {
	// 	return createUrlResult, err
	// }
	// if err := database.TxNamedQueryStruct(tx, query, url, &createUrlResult); err != nil {
	// 	if txErr := tx.Rollback(); txErr != nil {
	// 		return createUrlResult, fmt.Errorf("failed to rollback transaction: %w. root cause: %w", txErr, err)
	// 	}
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

// Get url by ID from db
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
			return DbUrl{}, ErrUrlNotFound
		}
		return DbUrl{}, err
	}
	return dbUrl, nil
}

// Delete url from DB
func (store Store) Delete(ctx context.Context, id uuid.UUID, txFunc TxFunc) error {
	queryData := struct {
		UrlId uuid.UUID `db:"url_id"`
	}{
		UrlId: id,
	}
	const query = `
		DELETE from urls
		WHERE url_id = :url_id
		RETURNING url_key
	`
	type UrlDeleteResult struct {
		Key string `db:"url_key"`
	}
	var urlDeleteRes UrlDeleteResult
	tx, err := database.WithTx(store.db)
	if err != nil {
		return err
	}
	if err := database.TxNamedQueryStruct(tx, query, queryData, &urlDeleteRes); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w. root cause: %w", txErr, err)
		}
		if errors.Is(err, database.ErrDbNotFound) {
			return ErrUrlNotFound
		}
		return err
	}
	if err = txFunc(urlDeleteRes.Key); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w. root cause: %w", txErr, err)
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
