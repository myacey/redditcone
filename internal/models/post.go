package models

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrCantMarshalPost = errors.New("cant marshal post")

var (
	postTypes      = []string{"text", "link"}
	postCategories = []string{"music", "funny", "videos", "programming", "news", "fashion"}
)

type Post struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Author    *User     `json:"author" bson:"author"`
	CreatedAt time.Time `json:"created" bson:"created"`

	Title    string `json:"title" bson:"title"`
	Category string `json:"category" bson:"category"`
	Text     string `json:"text,omitempty" bson:"text,omitempty"`
	URL      string `json:"url,omitempty" bson:"url,omitempty"`

	Views            int     `json:"views" bson:"views"`
	Score            int     `json:"score" bson:"score"`
	Type             string  `json:"type" bson:"type"`
	Votes            []*Vote `json:"votes" votes:"votes"`
	UpvotePercentage int     `json:"upvotePercentage" bson:"upvote_percantage"`

	CommentCount int        `json:"-" bson:"comment_count"`
	Comments     []*Comment `json:"comments" bson:"-"`
}

func NewPost(user *User, category, title, postType, postText, postURL string) *Post {
	newVote := NewVote(user.ID, 1)

	idWithHyphens := uuid.New().String()
	id := strings.ReplaceAll(idWithHyphens, "-", "")

	return &Post{
		ID:        id,
		Author:    user,
		CreatedAt: time.Now(),

		Title:    title,
		Category: category,

		Text: postText,
		URL:  postURL,

		Views:            0,
		Score:            1,
		Type:             postType,
		Votes:            []*Vote{newVote},
		UpvotePercentage: 100,

		Comments:     []*Comment{},
		CommentCount: 0,
	}
}

func (p *Post) GetMarshal() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}, ErrCantMarshalPost
	}

	return data, nil
}

// ValidatePost checks post type and post category
func ValidatePost(newPost Post) bool {
	return slices.Contains(postCategories, newPost.Category) && slices.Contains(postTypes, newPost.Type)
}

func GetCategories() []string {
	return postCategories
}

func AddNilComments(posts ...*Post) {
	for i := range posts {
		posts[i].Comments = make([]*Comment, posts[i].CommentCount)
		for j := range posts[i].Comments {
			posts[i].Comments[j] = &Comment{}
		}
	}
}
