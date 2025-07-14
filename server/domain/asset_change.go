package domain

import "time"

const (
	AssetChangeIDSuffix = "AssetChange"
)

type AssetChangeID string

func NewAssetChangeID() AssetChangeID {
	return AssetChangeID(NewUUIDv4(AssetChangeIDSuffix))
}

type AssetChange struct {
	ID        AssetChangeID
	UserID    UserID
	RecordID  RecordID
	AssetID   AssetID
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAssetChange(userID UserID, recordID RecordID, assetID AssetID, amount int) *AssetChange {
	return &AssetChange{
		ID:        NewAssetChangeID(),
		UserID:    userID,
		RecordID:  recordID,
		AssetID:   assetID,
		Amount:    amount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
