// Authentication manager for creating new tokens
package auth

import (
	"errors"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrEmptySecret = errors.New("cannot find jwt secret key in env variables")

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

// Create new Authentication manager with specified JWT secret and authentication time
func New(secret string, authTime int) Auth {
	auth := Auth{
		Secret:         secret,
		jwtAuthTimeSec: authTime,
	}
	return auth
}

// Create new token with specified claims
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
