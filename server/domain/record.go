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
	At          time.Time // 入出金が発生した日時（ユーザー指定）
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func newRecord(userID UserID, recordType RecordType, title, description string, at time.Time) *Record {
	return &Record{
		ID:          NewRecordID(),
		UserID:      userID,
		RecordType:  recordType,
		Title:       title,
		Description: description,
		At:          at,
	}
}

func NewRecordIncomeWithAssetChange(userID UserID, title string, description string, at time.Time, assetID AssetID, amount int) (*Record, *AssetChange, error) {
	if amount < 0 {
		return nil, nil, ErrInvalidRecordAmount
	}

	record := newRecord(userID, RecordTypeIncome, title, description, at)
	assetChange := NewAssetChange(userID, record.ID, assetID, amount)
	return record, assetChange, nil
}

func NewRecordExpenseWithAssetChange(userID UserID, title string, description string, at time.Time, assetID AssetID, amount int) (*Record, *AssetChange, error) {
	if amount < 0 {
		return nil, nil, ErrInvalidRecordAmount
	}

	record := newRecord(userID, RecordTypeExpense, title, description, at)
	assetChange := NewAssetChange(userID, record.ID, assetID, -amount)
	return record, assetChange, nil
}

func NewRecordTransferWithAssetChanges(userID UserID, title string, description string, at time.Time, fromAssetID AssetID, toAssetID AssetID, amount int) (*Record, *AssetChange, *AssetChange, error) {
	if amount < 0 {
		return nil, nil, nil, ErrInvalidRecordAmount
	}

	record := newRecord(userID, RecordTypeTransfer, title, description, at)

	fromAssetChange := NewAssetChange(userID, record.ID, fromAssetID, -amount)
	toAssetChange := NewAssetChange(userID, record.ID, toAssetID, amount)

	return record, fromAssetChange, toAssetChange, nil
}

type Records []*Record

func (records Records) OldestRecord(isReverse bool) (*Record, error) {
	if len(records) == 0 {
		return nil, ErrEntityNotFound
	}

	if isReverse {
		return records[len(records)-1], nil
	}
	return records[0], nil
}
