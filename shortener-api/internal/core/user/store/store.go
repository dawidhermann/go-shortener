package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
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
	ErrUniqueValidation     = errors.New("unique validation")
	ErrUniqueEmailViolation = errors.New("unique email violation")
	ErrUserNotFound         = errors.New("user not found")
)

type Store struct {
	db sqlx.ExtContext
}

type CreateUserResult struct {
	UserId    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store Store) Create(ctx context.Context, user DbUser) (CreateUserResult, error) {
	const query = `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3)
		RETURNING user_id, created_at
	`
	var createUserResult CreateUserResult
	if err := database.NamedQueryStruct(ctx, store.db, query, user, createUserResult); err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			switch pgerr.Code {
			case undefinedTable:
				return CreateUserResult{}, ErrUndefinedTable
			case uniqueViolation:
				return CreateUserResult{}, ErrUniqueValidation
			}
		}
		return createUserResult, err
	}
	return createUserResult, nil
}

func (store Store) Update(ctx context.Context, user DbUser) error {
	const query = `
		UPDATE users SET
			"email" = :email,
			"password" = :password
		WHERE
			user_id = :user_id
	`
	if err := database.NamedExecContext(ctx, store.db, query, user); err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == uniqueViolation {
				return ErrUniqueEmailViolation
			}
			return err
		}
		if err == database.ErrDoNoRowsAffected {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (store Store) GetById(ctx context.Context, id uuid.UUID) (DbUser, error) {
	queryData := struct {
		UserId uuid.UUID `db:"user_id"`
	}{
		UserId: id,
	}
	const query = `
		SELECT
			"user_id",
			"username",
			"email",
			"password",
			"enabled",
			"created_at"
		FROM users
		WHERE user_id = :user_id
	`
	var usr DbUser
	if err := database.NamedQueryStruct(ctx, store.db, query, queryData, &usr); err != nil {
		return DbUser{}, err
	}
	return usr, nil
}

func (store Store) GetByEmail(ctx context.Context, email mail.Address) (DbUser, error) {
	queryData := struct {
		Email mail.Address `db:"email"`
	}{
		Email: email,
	}
	const query = `
		SELECT
			"user_id",
			"username",
			"email",
			"password",
			"enabled",
			"date_created",
			"date_updated"
		FROM users
		WHERE email = :email
	`
	var usr DbUser
	if err := store.db.QueryRowxContext(ctx, query, queryData).Scan(&usr); err != nil {
		if err == sql.ErrNoRows {
			return DbUser{}, ErrUserNotFound
		}
		return DbUser{}, fmt.Errorf("failed to fetch user: %w", err)
	}
	return usr, nil
}

func (store Store) Delete(ctx context.Context, id uuid.UUID) error {
	queryData := struct {
		UserId uuid.UUID `db:"user_id"`
	}{
		UserId: id,
	}
	query := `
		DELETE FROM users
		WHERE "user_id" = :user_id
	`
	if err := database.NamedExecContext(ctx, store.db, query, queryData); err != nil {
		if err == database.ErrDoNoRowsAffected {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
