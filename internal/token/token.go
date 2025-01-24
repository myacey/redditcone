package token

import (
	"errors"

	"github.com/myacey/redditclone/internal/models"
)

var (
	ErrCantCreateToken    = errors.New("cant create token")
	ErrInvalidTokenMethod = errors.New("invalid token method")
	ErrTokenInvalid       = errors.New("invalid token")
	ErrNoClaims           = errors.New("no token claims")
	ErrTokenExpired       = errors.New("token expired")
	ErrCantExtractExpTime = errors.New("invalid token expiration time")
)

type TokenMaker interface {
	CreateToken(usr *models.User) (string, error)
	ExtractUserID(tokenString string) (string, error)
}
