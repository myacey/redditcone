package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/myacey/redditclone/internal/models"
)

var (
	ErrCommentAlreadyExists = errors.New("comment already exists")
	ErrCommentDontExists    = errors.New("comment dont exists")
)

func (s *Service) increatePostCommentCount(ctx context.Context, post *models.Post) error {
	post.CommentCount++

	return s.postRepo.UpdatePostInfo(ctx, post)
}

// createComment creates new comment
// intaraction with comment's functions excluesevly from posts
func (s *Service) createComment(ctx context.Context, newComment *models.Comment) error {
	// comment already exists
	_, err := s.commentRepo.GetCommentByID(ctx, newComment.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	return s.commentRepo.CreateComment(ctx, newComment)
}

// deleteComment deletes comment
// intaraction with comment's functions excluesevly from posts
func (s *Service) deleteComment(ctx context.Context, commentID string) error {
	return s.commentRepo.DeleteComment(ctx, commentID)
}

func (s *Service) AddCommentToPost(ctx context.Context, postID string, newComment models.Comment) (*models.Post, error) {
	if len(newComment.Body) == 0 {
		return nil, ErrCommentCantBeNull
	}

	// check if post exist
	gotPost, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// create comment with internal func
	err = s.createComment(ctx, &newComment)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	err = s.increatePostCommentCount(ctx, gotPost)
	if err != nil {
		return nil, err
	}

	gotPost.Comments = comments

	return gotPost, nil
}

func (s *Service) RemoveComment(ctx context.Context, postID, commentID string) (*models.Post, error) {
	// check if post really exists
	gotPost, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	err = s.deleteComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	gotPost.Comments = comments

	return gotPost, nil
}
