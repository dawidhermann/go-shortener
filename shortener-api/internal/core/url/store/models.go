package store

import (
	"time"

	"github.com/google/uuid"
)

type DbUrl struct {
	ID          uuid.UUID `db:"url_id"`
	Key         string    `db:"url_key"`
	UserId      uuid.UUID `db:"user_id"`
	DateCreated time.Time `db:"created_at"`
	DateUpdated time.Time `db:"updated_at"`
}
