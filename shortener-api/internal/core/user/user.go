package user

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/dawidhermann/shortener-api/internal/core/user/store"
	"github.com/dawidhermann/shortener-api/internal/sys/validate"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Core struct {
	store *store.Store
}

// Create new user code
func NewUserCore(db *sqlx.DB) *Core {
	return &Core{
		store: store.NewUserStore(db),
	}
}

// Validate user's data and create new user
func (core *Core) Create(ctx context.Context, userCreateModel UserCreateViewModel) (User, error) {
	err := validate.ValidateStruct(userCreateModel)
	if err != nil {
		return User{}, fmt.Errorf("user data validation error: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(userCreateModel.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	userData := User{
		Username: userCreateModel.Username,
		Password: hash,
		Email:    userCreateModel.Email,
	}
	result, err := core.store.Create(ctx, toDbUser(userData))
	if err != nil {
		return User{}, fmt.Errorf("store error: %w", err)
	}
	userData.UserId = result.UserId
	userData.DateCreated = result.DateCreated
	userData.DateUpdated = result.DateUpdated
	return userData, nil
}

// Validate user's data and updates existing user
func (core *Core) Update(ctx context.Context, userData UserPatchModel) error {
	err := validate.ValidateStruct(userData)
	user := User{}
	if err != nil {
		return fmt.Errorf("user data validation error: %w", err)
	}
	if userData.Email != nil {
		addr, err := mail.ParseAddress(*userData.Email)
		if err != nil {
			return fmt.Errorf("failed to parse email address: %w", err)
		}
		user.Email = *addr
	}
	if userData.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*userData.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hash
	}
	if err = core.store.Update(ctx, toDbUser(user)); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Fetch user by user's email
func (core *Core) GetByEmail(ctx context.Context, email mail.Address) (User, error) {
	userData, err := core.store.GetByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user by email: %w", err)
	}
	return toUser(userData), nil
}

// Fetch user by user's ID
func (core *Core) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	userData, err := core.store.GetById(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user by id: %w", err)
	}
	return toUser(userData), err
}

// Delete user by user's ID
func (core *Core) DeleteById(ctx context.Context, id uuid.UUID) error {
	if err := core.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user by id: %w", err)
	}
	return nil
}
