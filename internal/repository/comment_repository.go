package repository

import (
	"context"

	"github.com/myacey/redditclone/internal/models"
)

type CommentRepository interface {
	GetCommentByID(ctx context.Context, commentID string) (*models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*models.Comment, error)
	CreateComment(ctx context.Context, newComment *models.Comment) error
	DeleteComment(ctx context.Context, commentID string) error
}
