package auth

import (
	"net/mail"

	"github.com/golang-jwt/jwt/v5"
	// echojwt "github.com/labstack/echo-jwt/v4"
	// "github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	// "net/http"

	"errors"
	"time"
)

var ErrEmptySecret = errors.New("cannot find jwt secret key in env variables")

type KeyLookup interface {
	PrivateKeyPem() string
	PublicKeyPem() string
}

type Auth struct {
	Secret         string
	jwtAuthTimeSec int
}

type UserClaims struct {
	UserId string
	Email  *mail.Address
}

type JwtCustomClaims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func New(secret string, authTime int) Auth {
	auth := Auth{
		Secret:         secret,
		jwtAuthTimeSec: authTime,
	}
	return auth
}

func (auth Auth) NewToken(tokenClaims UserClaims) (string, error) {
	claims := &JwtCustomClaims{
		tokenClaims.UserId,
		tokenClaims.Email.String(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(auth.jwtAuthTimeSec))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(auth.Secret))
	if err != nil {
		return "", err
	}
	return t, nil
}
