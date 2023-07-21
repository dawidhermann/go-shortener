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

type KeyLookup interface {
	PrivateKeyPem() string
	PublicKeyPem() string
}

type Auth struct {
	keyLookup      KeyLookup
	jwtAuthTimeSec int
}

type TokenClaims struct {
	UserId string
	Email  *mail.Address
}

type jwtCustomClaims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func New(keyLookup KeyLookup, authTime int) Auth {
	auth := Auth{
		keyLookup:      keyLookup,
		jwtAuthTimeSec: authTime,
	}
	return auth
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

func (auth Auth) NewToken(tokenClaims TokenClaims) (string, error) {
	claims := &jwtCustomClaims{
		tokenClaims.UserId,
		tokenClaims.Email.String(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(auth.jwtAuthTimeSec))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	privateKey := auth.keyLookup.PrivateKeyPem()
	t, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

// func EncodeJwtToken(claims map[string]interface{}) (jwt.Token, string, error) {
// 	claimsMap := make(map[string]interface{})
// 	for key, value := range claims {
// 		claimsMap[key] = value
// 	}
// 	claimsMap["exp"] = time.Now().Add(authManager.tokenExpTime)
// 	return authManager.TokenAuth.Encode(claimsMap)
// }
