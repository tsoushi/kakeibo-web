package domain

import (
	"time"
)

const (
	RecordIDSuffix = "User"
)

type RecordID string

func NewRecordID() RecordID {
	return RecordID(NewUUIDv4(RecordIDSuffix))
}

type RecordType string

const (
	RecordTypeExpense  RecordType = "EXPENSE"
	RecordTypeIncome   RecordType = "INCOME"
	RecordTypeTransfer RecordType = "TRANSFER"
)

type Record struct {
	ID          RecordID
	UserID      UserID
	RecordType  RecordType
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func newRecord(userID UserID, recordType RecordType, title, description string) *Record {
	return &Record{
		ID:          NewRecordID(),
		UserID:      userID,
		RecordType:  recordType,
		Title:       title,
		Description: description,
	}
}

func NewRecordExpenseOrIncomeWithAssetChange(userID UserID, title string, description string, assetID AssetID, amount int) (*Record, *AssetChange) {
	var recordType RecordType = RecordTypeIncome
	if amount < 0 {
		recordType = RecordTypeExpense
	}

	record := newRecord(userID, recordType, title, description)

	assetChange := NewAssetChange(userID, record.ID, assetID, amount)

	return record, assetChange
}

func NewRecordTransferWithAssetChanges(userID UserID, title string, description string, fromAssetID AssetID, toAssetID AssetID, amount int) (*Record, *AssetChange, *AssetChange, error) {
	if amount < 0 {
		return nil, nil, nil, ErrInvalidRecordAmount
	}

	record := newRecord(userID, RecordTypeTransfer, title, description)

	fromAssetChange := NewAssetChange(userID, record.ID, fromAssetID, -amount)
	toAssetChange := NewAssetChange(userID, record.ID, toAssetID, amount)

	return record, fromAssetChange, toAssetChange, nil
}
