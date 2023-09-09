// Definition of endpoints for authentcating users
package authctrl

import (
	"errors"
	"net/http"
	"net/mail"

	v1 "github.com/dawidhermann/shortener-api/api/v1"
	"github.com/dawidhermann/shortener-api/internal/auth"
	"github.com/dawidhermann/shortener-api/internal/core/user"
	"github.com/dawidhermann/shortener-api/internal/core/user/store"
	"github.com/dawidhermann/shortener-api/internal/sys/validate"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrBasicAuthParse      = errors.New("failed to parse basic auth")
	ErrLoginData           = errors.New("incorrect login data")
	ErrPasswordNotMatch    = errors.New("incorrect password")
	ErrAuthenticationError = errors.New("authentication failed")
)

type AuthController struct {
	Auth auth.Auth
	Core *user.Core
}

// Returns JWT token if user exists and credential match
func (ctrl AuthController) LoginUser(c echo.Context) error {
	email, password, ok := c.Request().BasicAuth()
	if !ok {
		return v1.NewRequestError(ErrBasicAuthParse, http.StatusBadRequest)
	}
	loginModel := struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}{
		Email:    email,
		Password: password,
	}
	err := validate.ValidateStruct(loginModel)
	if err != nil {
		return v1.NewRequestError(ErrLoginData, http.StatusBadRequest)
	}
	userEmail, err := mail.ParseAddress(loginModel.Email)
	if err != nil {
		return v1.NewRequestError(ErrLoginData, http.StatusBadRequest)
	}
	user, err := ctrl.Core.GetByEmail(c.Request().Context(), *userEmail)
	if errors.Is(err, store.ErrUserNotFound) {
		return v1.NewRequestError(err, http.StatusNotFound)
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return v1.NewRequestError(ErrPasswordNotMatch, http.StatusBadRequest)
	}
	token, err := ctrl.Auth.NewToken(auth.UserClaims{
		UserId: user.UserId.String(),
		Email:  &user.Email,
	})
	if err != nil {
		return v1.NewRequestError(ErrAuthenticationError, http.StatusInternalServerError)
	}
	tokenResp := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	return c.JSON(http.StatusOK, tokenResp)
}
