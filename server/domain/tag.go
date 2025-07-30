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

type Tags []*Tag

func (t Tags) ContainsByName(name string) bool {
	for _, tag := range t {
		if tag.Name == name {
			return true
		}
	}
	return false
}

func NewTagsNotExist(userID UserID, existTags Tags, names []string) []*Tag {
	tags := make([]*Tag, 0, len(names))
	for _, name := range names {
		if !existTags.ContainsByName(name) {
			tags = append(tags, NewTag(userID, name))
		}
	}

	return tags
}

type TagWithRecordID struct {
	Tag
	RecordID RecordID
}
