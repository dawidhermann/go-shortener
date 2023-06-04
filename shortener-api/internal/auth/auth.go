package auth

import (
	"database/sql"
	"errors"
	"github.com/dawidhermann/shortener-api/config"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type AuthManager struct {
	dbConn *sql.DB
}

type DbConfig struct {
	DbUser     string
	DbPassword string
	DbAddr     string
	DbName     string
}

//func NewAuthManager(db *sql.DB) *AuthManager {
//	return &AuthManager{dbConn: db}
//}
//
//func (authMng *AuthManager) GetUserByUsername(username string) error {
//
//}

//func (authMng *AuthManager) Close() {
//	authMng.Close()
//}

var ErrEmptySecret = errors.New("cannot find jwt secret key in env variables")

func NewAuthenticationManager(authConfig config.AuthConfig) (AuthenticationManager, error) {
	if len(authConfig.JwtSecretKey) == 0 {
		return AuthenticationManager{}, ErrEmptySecret
	}
	tokenAuth := jwtauth.New("HS256", []byte(authConfig.JwtSecretKey), nil)
	return AuthenticationManager{TokenAuth: tokenAuth, tokenExpTime: time.Duration(authConfig.JwtExpTime) * time.Second}, nil
}

func (authManager AuthenticationManager) EncodeJwtToken(claims map[string]interface{}) (jwt.Token, string, error) {
	claimsMap := make(map[string]interface{})
	for key, value := range claims {
		claimsMap[key] = value
	}
	claimsMap["exp"] = time.Now().Add(authManager.tokenExpTime)
	return authManager.TokenAuth.Encode(claimsMap)
}
