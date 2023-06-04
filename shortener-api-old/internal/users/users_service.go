package users

import (
	"errors"
	"github.com/dawidhermann/shortener-api/internal/db"
	"github.com/dawidhermann/shortener-api/internal/validators"
)

var (
	ErrPasswordMismatch = errors.New("confirmation password does not match password")
	ErrInvalidPassword  = errors.New("password is not valid")
	ErrInvalidUsername  = errors.New("username is not valid")
	ErrInvalidEmail     = errors.New("email is not valid")
	ErrIncorrectUserId  = errors.New("incorrect user id")
	ErrUserNotFound     = errors.New("user not found")
)

type ServiceUsers struct {
	repository RepositoryUsers
}

func NewServiceUsers(connDb db.SqlConnection) ServiceUsers {
	return ServiceUsers{
		repository: newRepositoryUsers(connDb),
	}
}

func (service ServiceUsers) CreateUser(createUserModel UserCreateViewModel) (int, error) {
	if !validators.ValidateUserPassword(createUserModel.Password) {
		return 0, ErrInvalidPassword
	}
	if createUserModel.Password != createUserModel.PasswordConfirm {
		return 0, ErrPasswordMismatch
	}
	if !validators.ValidateUserName(createUserModel.Username) {
		return 0, ErrInvalidUsername
	}
	if !validators.ValidateEmail(createUserModel.Email) {
		return 0, ErrInvalidEmail
	}
	userId, err := service.repository.createUserEntity(createUserModel.Username, createUserModel.Password, createUserModel.Email)
	return userId, err
}

func (service ServiceUsers) GetUser(userId int) (user, error) {
	userData, err := service.repository.getUserEntity(userId)
	if err != nil {
		return user{}, ErrUserNotFound
	}
	return userData, nil
}

func (service ServiceUsers) GetUserByUsername(username string) (user, error) {
	userData, err := service.repository.getUserEntityByUsername(username)
	if err != nil {
		return user{}, ErrUserNotFound
	}
	return userData, nil
}

func (service ServiceUsers) DeleteUser(userId int) error {
	if userId == 0 {
		return ErrIncorrectUserId
	}
	return service.repository.deleteUserEntity(userId)
}

func (service ServiceUsers) UpdateUser(userId int, userPatchModel UserPatchModel) error {
	userData, err := service.GetUser(userId)
	if err != nil {
		return err
	}
	if userPatchModel.Email != nil {
		email := *userPatchModel.Email
		if !validators.ValidateEmail(email) {
			return ErrInvalidEmail
		}
		userData.Email = email
	}
	var patchModelPassword string
	var patchModelPasswordConfirmation string
	if userPatchModel.Password != nil {
		patchModelPassword = *userPatchModel.Password
		if !validators.ValidateUserPassword(patchModelPassword) {
			return ErrInvalidPassword
		}
		if userPatchModel.PasswordConfirm != nil {
			patchModelPasswordConfirmation = *userPatchModel.PasswordConfirm
		}
		if patchModelPassword != patchModelPasswordConfirmation {
			return ErrPasswordMismatch
		}
		userData.Password = patchModelPassword
	} else if userPatchModel.PasswordConfirm != nil {
		return ErrPasswordMismatch
	}
	err = service.repository.updateUserEntity(userData)
	return err
}
