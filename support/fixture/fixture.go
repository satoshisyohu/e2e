// Package fixture はテストの前提データを公開 API 経由で作成する。
//
// 方針:
//   - データはすべて公開 API で作る(DB への直接投入はしない)。
//   - 値には idgen を使い、他のテストや他の実行と衝突させない。
//   - 作成に失敗したら t.Fatal で即座にテストを止める(前提が崩れているため)。
//   - 削除は行わない。データは idgen の実行IDで分離されている。
package fixture

import (
	"context"

	"e2e/support/client"
)

// Fixture はテストデータを作成するヘルパー。
type Fixture struct {
	clients *client.Clients
}

// New は Fixture を生成する。
func New(c *client.Clients) *Fixture {
	return &Fixture{clients: c}
}

func ctx() context.Context { return context.Background() }
