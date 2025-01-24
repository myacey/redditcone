package service

import (
	"context"
	"net/http"
	"time"

	"github.com/myacey/redditclone/internal/customerror/errhandler"
	"github.com/myacey/redditclone/internal/models"
)

// func (s *Service) AddUserToDB(ctx context.Context, user *models.User) error {
// 	return s.userRepo.CreateUser(ctx, user)
// }

func (s *Service) GetUserFromDBByID(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *Service) GetUserFromDBByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(ctx, username)
}

func (s *Service) CreateNewUser(ctx context.Context, user *models.User) (*models.Session, error) {
	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, errhandler.New(http.StatusInternalServerError, "internal error", "cant create user in db", err)
	}

	token, err := s.tokenMaker.CreateToken(user)
	if err != nil {
		return nil, errhandler.New(http.StatusInternalServerError, "internal error", "cant create token", err)
	}
	s.logger.Infow("created token",
		"username", user.Username,
		"user_id", user.ID,
		"exp", time.Now().Add(time.Hour*24),
	)

	session := models.NewSession(token)
	err = s.sessionRepo.CreateSession(ctx, session, user.ID, 24*time.Hour)
	if err != nil {
		return nil, errhandler.New(http.StatusInternalServerError, "internal error", "cant create token in db", err)
	}

	return session, nil
}

func (s *Service) LoginUser(ctx context.Context, username string) (*models.Session, error) {
	usr, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenMaker.CreateToken(usr)
	if err != nil {
		return nil, err
	}
	s.logger.Infow("created token",
		"username", usr.Username,
		"user_id", usr.ID,
		"exp", time.Now().Add(time.Hour*24),
	)

	session := models.NewSession(token)

	err = s.sessionRepo.UpdateSessionToken(ctx, session, usr.ID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return session, nil
}
