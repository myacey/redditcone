package service

import (
	"context"
	"fmt"
)

func (s *Service) CheckUserSession(ctx context.Context, userID, token string) error {
	dbToken, err := s.sessionRepo.GetSessionTokenByUsername(ctx, userID)
	if err != nil {
		return err
	}

	if dbToken == token {
		return nil
	}

	return fmt.Errorf("invalid token")
}
