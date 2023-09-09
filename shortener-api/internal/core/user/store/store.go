package store

import (
	"context"
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
	ErrUniqueViolation      = errors.New("unique violation")
	ErrUniqueEmailViolation = errors.New("unique email violation")
	ErrUserNotFound         = errors.New("user not found")
)

type Store struct {
	db sqlx.ExtContext
}

type CreateUserResult struct {
	UserId    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt time.Time
}

func NewUserStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store Store) Create(ctx context.Context, user DbUser) (CreateUserResult, error) {
	const query = `
		INSERT INTO users (username, password, email)
		VALUES (:username, :password, :email)
		RETURNING user_id, created_at
	`
	var createUserResult CreateUserResult
	if err := database.NamedQueryStruct(ctx, store.db, query, user, &createUserResult); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case undefinedTable:
				return CreateUserResult{}, ErrUndefinedTable
			case uniqueViolation:
				return CreateUserResult{}, ErrUniqueViolation
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
		if err == database.ErrDbNoRowsAffected {
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
		if errors.Is(err, database.ErrDbNotFound) {
			return DbUser{}, ErrUserNotFound
		}
		return DbUser{}, err
	}
	return usr, nil
}

func (store Store) GetByEmail(ctx context.Context, email mail.Address) (DbUser, error) {
	queryData := struct {
		Email string `db:"email"`
	}{
		Email: email.Address,
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
		WHERE email = :email
	`
	var usr DbUser
	if err := database.NamedQueryStruct(ctx, store.db, query, queryData, &usr); err != nil {
		fmt.Println(err)
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
	const query = `
		DELETE FROM users
		WHERE user_id = :user_id
	`
	if err := database.NamedExecContext(ctx, store.db, query, queryData); err != nil {
		if err == database.ErrDbNoRowsAffected {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
