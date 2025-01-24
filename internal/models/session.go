package models

import (
	"encoding/json"
	"errors"
)

var (
	ErrCantMarshal = errors.New("cant marshal token")
)

type Session struct {
	Token string `json:"token" redis:"token"`
}

func NewSession(token string) *Session {
	return &Session{token}
}

func (s *Session) GetMarshal() ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return []byte{}, ErrCantMarshal
	}

	return data, nil
}
