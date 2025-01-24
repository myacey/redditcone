package service

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/myacey/redditclone/internal/customerror/errhandler"
	"github.com/myacey/redditclone/internal/models"
)

var (
	ErrInvalidPostData   = errors.New("invalid post data")
	ErrUnknown           = errors.New("unknown error")
	ErrCommentCantBeNull = errors.New("comment cant be null")
)

func (s *Service) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	posts, err := s.postRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}
	models.AddNilComments(posts...) // to show comment count
	return posts, nil
}

func (s *Service) AddPost(ctx context.Context, newPost *models.Post) error {
	if !models.ValidatePost(*newPost) {
		return ErrInvalidPostData
	}

	return s.postRepo.CreatePost(ctx, newPost)
}

func (s *Service) GetPostByID(ctx context.Context, postID string, increateVote bool) (*models.Post, error) {
	gotPost, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, errhandler.New(http.StatusBadRequest, "cant find post", "invalid params to find post", err)
	}

	if increateVote {
		err = s.increatePostViews(ctx, gotPost)
		if err != nil {
			return nil, err
		}
	}

	comments, err := s.commentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, errhandler.New(http.StatusBadRequest, "cant find comments", "invalid params to find comments", err)
	}

	gotPost.Comments = comments

	return gotPost, nil
}

func (s *Service) increatePostViews(ctx context.Context, post *models.Post) error {
	post.Views++

	return s.postRepo.UpdatePostInfo(ctx, post)
}

func (s *Service) GetPostsByAuthor(ctx context.Context, username string) ([]*models.Post, error) {
	// check if user really exists
	_, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	allPosts, err := s.postRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}
	sortedPosts := []*models.Post{}
	for _, v := range allPosts {
		if v.Author.Username == username {
			sortedPosts = append(sortedPosts, v)
		}
	}
	models.AddNilComments(sortedPosts...) // to show comment count

	return sortedPosts, nil
}

func (s *Service) GetPostsByCategory(ctx context.Context, category string) ([]*models.Post, error) {
	if !slices.Contains(models.GetCategories(), category) {
		return nil, errhandler.New(http.StatusBadRequest, "invalid category", "invalid category", nil)
	}

	allPosts, err := s.postRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}
	sortedPosts := []*models.Post{}
	for _, v := range allPosts {
		if v.Category == category {
			sortedPosts = append(sortedPosts, v)
		}
	}
	models.AddNilComments(sortedPosts...) // to show comment count

	return sortedPosts, nil
}

func (s *Service) DeletePostWithID(ctx context.Context, postID string) error {
	return s.postRepo.DeletePost(ctx, postID)
}
