package repository

import (
	"context"
	"time"

	"github.com/myacey/redditclone/internal/models"
)

type SessionRepository interface {
	CreateSession(
		ctx context.Context,
		session *models.Session,
		username string,
		expirationTime time.Duration,
	) error
	GetSessionTokenByUsername(ctx context.Context, username string) (string, error)
	UpdateSessionToken(
		ctx context.Context,
		newSession *models.Session,
		username string,
		expirationTime time.Duration,
	) error
}
