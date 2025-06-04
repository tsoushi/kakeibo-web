package domain

import "time"

const (
	AssetIDSuffix = "Asset"
)

type AssetID string

func NewAssetID() AssetID {
	return AssetID(NewUUIDv4(AssetIDSuffix))
}

type Asset struct {
	ID         AssetID
	UserID     UserID
	Name       string
	CategoryID *AssetCategoryID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewAsset(userID UserID, name string, categoryID *AssetCategoryID) *Asset {
	return &Asset{
		ID:         NewAssetID(),
		UserID:     userID,
		Name:       name,
		CategoryID: categoryID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
