package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/myacey/redditclone/internal/models"
	"github.com/myacey/redditclone/internal/repository"
)

type RedisSessionRepo struct {
	rdb *redis.Client
}

func NewRedisSessionRepo(rdb *redis.Client) repository.SessionRepository {
	return &RedisSessionRepo{rdb: rdb}
}

func (r *RedisSessionRepo) CreateSession(
	ctx context.Context,
	session *models.Session,
	userID string,
	expirationTime time.Duration,
) error {
	marshalled, err := session.GetMarshal()
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, userID, marshalled, expirationTime).Err()
}

func (r *RedisSessionRepo) GetSessionTokenByUsername(ctx context.Context, userID string) (string, error) {
	var session *models.Session
	marshalledSession, err := r.rdb.Get(ctx, userID).Result()
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal([]byte(marshalledSession), &session); err != nil {
		return "", err
	}

	return session.Token, nil
}

func (r *RedisSessionRepo) UpdateSessionToken(
	ctx context.Context,
	newSession *models.Session,
	userID string,
	expirationTime time.Duration,
) error {
	marshalled, err := json.Marshal(newSession)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, userID, marshalled, expirationTime).Err()
}
