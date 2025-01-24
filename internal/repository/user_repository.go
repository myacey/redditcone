package repository

import (
	"context"
	"errors"

	"github.com/myacey/redditclone/internal/models"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists in db")
	ErrUserDontExists    = errors.New("user dont exist in db")
	ErrInvalidToken      = errors.New("invalid token")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}
