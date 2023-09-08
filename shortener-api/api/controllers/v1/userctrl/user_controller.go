package userctrl

import (
	"errors"
	"fmt"
	"net/http"

	v1 "github.com/dawidhermann/shortener-api/api/v1"
	"github.com/dawidhermann/shortener-api/internal/core/user"
	"github.com/dawidhermann/shortener-api/internal/core/user/store"

	// "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrBindingRequestData = errors.New("failed to bind request data")
	ErrInvalidId          = errors.New("invalid id")
)

type UsersController struct {
	Core *user.Core
}

func (ctrl UsersController) CreateUser(c echo.Context) error {
	var userCreateModel user.UserCreateViewModel
	if err := c.Bind(&userCreateModel); err != nil {
		return v1.NewRequestError(ErrBindingRequestData, http.StatusBadRequest)
	}
	user, err := ctrl.Core.Create(c.Request().Context(), userCreateModel)
	if err != nil {
		if errors.Is(err, store.ErrUniqueViolation) {
			return v1.NewRequestError(store.ErrUniqueViolation, http.StatusConflict)
		}
		return fmt.Errorf("failed to create user [%v], %w", &user, err)
	}
	return c.JSON(http.StatusCreated, user)
}

func (controller UsersController) UpdateUser(c echo.Context) error {
	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return v1.NewRequestError(ErrInvalidId, http.StatusBadRequest)
	}
	var userPatchModel user.UserPatchModel
	if err := c.Bind(&userPatchModel); err != nil {
		return v1.NewRequestError(ErrBindingRequestData, http.StatusBadRequest)
	}
	if err := controller.Core.Update(c.Request().Context(), userId, userPatchModel); err != nil {
		if err == user.ErrUserNotValid {
			return v1.NewRequestError(user.ErrUserNotValid, http.StatusBadRequest)
		}
		if errors.Is(err, store.ErrUserNotFound) {
			return v1.NewRequestError(store.ErrUserNotFound, http.StatusNotFound)
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	userLocation := fmt.Sprintf("/users/%s", userId)
	responseHeader := c.Response().Header()
	responseHeader.Set("Location", userLocation)
	responseHeader.Set("Content-Location", userLocation)
	return c.NoContent(http.StatusNoContent)
}

func (controller UsersController) GetUserById(c echo.Context) error {
	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return v1.NewRequestError(ErrInvalidId, http.StatusBadRequest)
	}
	userData, err := controller.Core.GetById(c.Request().Context(), userId)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return v1.NewRequestError(store.ErrUserNotFound, http.StatusNotFound)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}
	userViewModel := user.NewUserViewModel(userData)
	return c.JSON(http.StatusOK, userViewModel)
}

func (controller UsersController) DeleteUser(c echo.Context) error {
	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return v1.NewRequestError(ErrInvalidId, http.StatusBadRequest)
	}
	err = controller.Core.DeleteById(c.Request().Context(), userId)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return v1.NewRequestError(store.ErrUserNotFound, http.StatusNotFound)
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return c.NoContent(http.StatusOK)
}
