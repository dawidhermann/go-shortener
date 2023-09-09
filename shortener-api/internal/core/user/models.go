package user

import (
	"net/mail"
	"time"

	"github.com/dawidhermann/shortener-api/internal/core/user/store"

	"github.com/google/uuid"
)

type UserViewModel struct {
	UserId   uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserCreateViewModel struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"eqfield=Password"`
}

type UserPatchModel struct {
	Email           *string `json:"email" validate:"omitempty,email"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"passwordConfirm" validate:"omitempty,eqfield=Password"`
}

type User struct {
	UserId      uuid.UUID
	Username    string
	Password    []byte
	Email       mail.Address
	Enabled     bool
	DateCreated time.Time
	DateUpdated time.Time
}

// Create new user view model
func NewUserViewModel(user User) UserViewModel {
	return UserViewModel{
		UserId:   user.UserId,
		Username: user.Username,
		Email:    user.Email.Address,
	}
}

func toDbUser(user User) store.DbUser {
	return store.DbUser{
		ID:          user.UserId,
		Username:    user.Username,
		Email:       user.Email.Address,
		Password:    user.Password,
		Enabled:     user.Enabled,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}
}

func toUser(userData store.DbUser) User {
	emailAddr := mail.Address{
		Address: userData.Email,
	}
	return User{
		UserId:      userData.ID,
		Username:    userData.Username,
		Password:    userData.Password,
		Email:       emailAddr,
		Enabled:     userData.Enabled,
		DateCreated: userData.DateCreated,
		DateUpdated: userData.DateUpdated,
	}
}
