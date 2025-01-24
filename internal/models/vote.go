package models

type Vote struct {
	UserID string `json:"user"`
	Vote   int8   `json:"vote"`
}

func NewVote(userID string, vote int8) *Vote {
	return &Vote{userID, vote}
}
