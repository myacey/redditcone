package service

import (
	"context"
	"errors"

	"github.com/myacey/redditclone/internal/models"
)

var (
	ErrVoteDontExist     = errors.New("vote dont exists")
	ErrVoteAlreadyExists = errors.New("vote already exists")
)

func (s *Service) updateVoteStat(post *models.Post) {
	if len(post.Votes) == 0 {
		post.Score = 0
		post.UpvotePercentage = 0
		return
	}
	positive := 0
	score := 0

	for _, v := range post.Votes {
		score += int(v.Vote)
		if v.Vote == 1 {
			positive++
		}
	}

	upvotePercentage := int((float32(positive) / float32(len(post.Votes)) * 100))

	post.UpvotePercentage = upvotePercentage
	post.Score = score
}

func (s *Service) VotePostWithID(ctx context.Context, postID string, newVote *models.Vote) (*models.Post, error) {
	gotPost, err := s.GetPostByID(ctx, postID, false)
	if err != nil {
		return nil, err
	}

	for i, v := range gotPost.Votes {
		if v.UserID != newVote.UserID {
			continue
		}

		if v.Vote == newVote.Vote { // post already exists
			return nil, ErrVoteAlreadyExists
		}

		gotPost.Votes[i] = newVote // we just need to change vote/upvote
		s.updateVoteStat(gotPost)
		err = s.postRepo.UpdatePostInfo(ctx, gotPost)
		if err != nil {
			return nil, err
		}

		return gotPost, nil
	}

	gotPost.Votes = append(gotPost.Votes, newVote) // create new vote

	s.updateVoteStat(gotPost)
	err = s.postRepo.UpdatePostInfo(ctx, gotPost)
	if err != nil {
		return nil, err
	}

	return gotPost, nil
}

func (s *Service) UnvotePostWithID(ctx context.Context, postID, userID string) (*models.Post, error) {
	gotPost, err := s.GetPostByID(ctx, postID, false)
	if err != nil {
		return nil, err
	}

	for i, v := range gotPost.Votes {
		if v.UserID == userID {
			gotPost.Votes = append(gotPost.Votes[:i], gotPost.Votes[i+1:]...) // delete vote

			s.updateVoteStat(gotPost) // new vote stat to post

			err := s.postRepo.UpdatePostInfo(ctx, gotPost) // apply changes to DB
			if err != nil {
				return nil, err
			}
			return gotPost, nil
		}
	}

	return nil, ErrVoteDontExist
}
