package store

import (
	"time"

	"github.com/google/uuid"
)

type DbUser struct {
	ID          uuid.UUID `db:"user_id"`
	Username    string    `db:"username"`
	Password    []byte    `db:"password"`
	Email       string    `db:"email"`
	Enabled     bool      `db:"enabled"`
	DateCreated time.Time `db:"created_at"`
	DateUpdated time.Time `db:"updated_at"`
}
