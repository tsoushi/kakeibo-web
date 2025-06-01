package domain

import "time"

type AssetCategoryID string

func NewAssetCategoryID() AssetCategoryID {
	return AssetCategoryID(NewUUIDv4("AssetCategory"))
}

type AssetCategory struct {
	ID        AssetCategoryID
	UserID    UserID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAssetCategory(userID UserID, name string) *AssetCategory {
	return &AssetCategory{
		ID:        NewAssetCategoryID(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
