package domain

import "time"

const (
	TagIDSuffix = "Tag"
)

type TagID string

func NewTagID() TagID {
	return TagID(NewUUIDv4(TagIDSuffix))
}

type Tag struct {
	ID        TagID
	UserID    UserID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTag(userID UserID, name string) *Tag {
	return &Tag{
		ID:        NewTagID(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
