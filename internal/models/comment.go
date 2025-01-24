package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Body      string    `json:"body" bson:"body"`
	Author    *User     `json:"author" bson:"author"`
	CreatedAt time.Time `json:"created" bson:"created"`

	RelatedPostID string `bson:"post_id"`
}

func NewComment(text string, author *User, postID string) *Comment {
	idWithHyphens := uuid.New().String()
	id := strings.ReplaceAll(idWithHyphens, "-", "")

	return &Comment{
		ID:            id,
		Body:          text,
		Author:        author,
		CreatedAt:     time.Now(),
		RelatedPostID: postID,
	}
}
