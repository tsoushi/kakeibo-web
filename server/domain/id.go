package domain

import (
	"fmt"
	"kakeibo-web-server/lib/id"
)

const (
	IDSeparator = "_"
)

type ID string

// 実行されるたびにランダムなUUIDを生成する
// TODO: 冪等性を持たせるためにUUIDv5を作成することも検討する(同じseedからは同じUUIDが生成されるような要件が必要な場合)
func NewUUIDv4(suffix string) ID {
	return ID(fmt.Sprintf("%s%s%s", id.NewUUIDv4(), IDSeparator, suffix))
}
