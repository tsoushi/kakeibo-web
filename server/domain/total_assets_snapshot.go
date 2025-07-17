package domain

import "time"

const (
	TotalAssetsSnapshotIDSuffix    = "TotalAssetsSnapshot"
	TotalAssetsSnapshotDefaultSpan = 5 // デフォルトのスナップショット間隔（件数）
)

type TotalAssetsSnapshotID string

func NewTotalAssetsSnapshotID() TotalAssetsSnapshotID {
	return TotalAssetsSnapshotID(NewUUIDv4(TotalAssetsSnapshotIDSuffix))
}

type TotalAssetsSnapshot struct {
	ID        TotalAssetsSnapshotID
	UserID    UserID
	AssetID   *AssetID
	At        time.Time // この時点での資産の合計額を記録
	Amount    int
	IsValid   bool // Snapshot作成後にAt以前のAssetChangeが変更された場合、この値はfalseになる
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTotalAssetsSnapshot(userID UserID, assetID *AssetID, at time.Time, amount int) *TotalAssetsSnapshot {
	return &TotalAssetsSnapshot{
		ID:        NewTotalAssetsSnapshotID(),
		UserID:    userID,
		AssetID:   assetID,
		At:        at,
		Amount:    amount,
		IsValid:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
