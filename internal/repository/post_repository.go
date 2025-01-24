package repository

import (
	"context"
	"errors"

	"github.com/myacey/redditclone/internal/models"
)

var (
	ErrPostAlreadyExists = errors.New("post already exists")
	ErrPostDontExists    = errors.New("post dont exists")

	ErrCommentAlreadyExists = errors.New("comment already exists")
	ErrCommentDontExists    = errors.New("comment dont exist")
)

type PostRepository interface {
	CreatePost(ctx context.Context, newPost *models.Post) error
	GetAllPosts(ctx context.Context) ([]*models.Post, error)
	GetPostByID(ctx context.Context, postID string) (*models.Post, error)
	UpdatePostInfo(ctx context.Context, updatedPost *models.Post) error
	DeletePost(ctx context.Context, postID string) error
}
