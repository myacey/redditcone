package models

import (
	"strings"

	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"`
}

func NewUser(username, password string) *User {
	idWithHyphens := uuid.New().String()
	id := strings.ReplaceAll(idWithHyphens, "-", "")
	return &User{
		ID:       id,
		Username: username,
		Password: password,
	}
}
