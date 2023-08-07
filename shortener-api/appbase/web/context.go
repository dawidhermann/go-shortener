package web

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/dawidhermann/shortener-api/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	ErrClaimNotFound      = errors.New("claim not found")
	ErrClaimTypeAssertion = errors.New("claim has different type")
)

func GetUserClaims(c echo.Context) (auth.UserClaims, error) {
	claims := getTokenClaims(c)
	userIdClaim, err := getStringClaim(claims, "userId")
	if err != nil {
		return auth.UserClaims{}, fmt.Errorf("failed to get userId claim: %w", ErrClaimNotFound)
	}
	emailClaim, err := getStringClaim(claims, "email")
	if err != nil {
		return auth.UserClaims{}, fmt.Errorf("failed to get email claim: %w", ErrClaimNotFound)
	}
	emailAddr, err := mail.ParseAddress(emailClaim)
	if err != nil {
		return auth.UserClaims{}, fmt.Errorf("failed to parse email from token claims: %w", err)
	}
	return auth.UserClaims{
		UserId: userIdClaim,
		Email:  emailAddr,
	}, nil
}

func getTokenClaims(c echo.Context) jwt.MapClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}

func getStringClaim(claims jwt.MapClaims, claimName string) (string, error) {
	if rawClaim, ok := claims[claimName]; ok {
		if claim, ok := rawClaim.(string); ok {
			return claim, nil
		}
		return "", ErrClaimTypeAssertion
	}
	return "", ErrClaimNotFound
}
