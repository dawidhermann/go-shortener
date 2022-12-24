package users

import (
	"errors"
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

type usersService struct {
}

func CreateUser(createUserModel UserCreateViewModel) (int, error) {
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
	userId, err := createUserEntity(createUserModel.Username, createUserModel.Password, createUserModel.Email)
	return userId, err
}

func GetUser(userId string) (user, error) {
	userData, err := getUserEntity(userId)
	if err != nil {
		return user{}, ErrUserNotFound
	}
	return userData, nil
}

func DeleteUser(userId string) error {
	if len(userId) == 0 {
		return ErrIncorrectUserId
	}
	return deleteUserEntity(userId)
}

func UpdateUser(userId string, userPatchModel UserPatchModel) error {
	userData, err := GetUser(userId)
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
	err = updateUserEntity(userData)
	return err
}
