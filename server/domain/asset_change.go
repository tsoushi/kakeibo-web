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

type AssetChanges []*AssetChange

type AssetChangeWithAt struct {
	AssetChange
	At time.Time // RecordのAtを持つ
}

type AssetChangeWithAts []*AssetChangeWithAt

func (changes AssetChangeWithAts) TotalAmount() int {
	total := 0
	for _, change := range changes {
		total += change.Amount
	}
	return total
}

func (changes AssetChangeWithAts) CreateSnapshots(assetID *AssetID, initAmount, span int) []*TotalAssetsSnapshot {
	if len(changes) == 0 {
		return nil
	}

	snapshots := make([]*TotalAssetsSnapshot, 0)
	currentAmount := initAmount

	count := 0
	lastAt := time.Time{}
	for _, change := range changes {
		// スナップショットのAt以前のAssetChangeをすべて含むことを保証するために直前のAssetChangeのAtと同値でないことを確認する
		if count >= span && change.At.After(lastAt) {
			snapshot := NewTotalAssetsSnapshot(change.UserID, assetID, lastAt, currentAmount)
			snapshots = append(snapshots, snapshot)
			count = 0
		}

		if assetID != nil && change.AssetID != *assetID {
			continue
		}
		count++

		currentAmount += change.Amount
		lastAt = change.At
	}

	return snapshots
}
