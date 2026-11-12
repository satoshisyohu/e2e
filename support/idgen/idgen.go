// Package idgen はテストデータを他の実行・他のテストと衝突させないための一意な値を生成する。
//
// 生成される値はすべて実行ID(RunID)を含むため、データがどの実行で作られたかを追跡できる。
package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

var runID = initRunID()

func initRunID() string {
	// CI では E2E_RUN_ID にパイプラインIDなどを渡すと追跡しやすい。
	if v := os.Getenv("E2E_RUN_ID"); v != "" {
		return v
	}
	return time.Now().Format("20060102150405") + "-" + randHex(2)
}

// RunID は今回のテスト実行を識別する ID を返す。
func RunID() string { return runID }

// New は "e2e-<runID>-<prefix>-<random>" 形式の一意な文字列を返す。
func New(prefix string) string {
	return fmt.Sprintf("e2e-%s-%s-%s", runID, prefix, randHex(4))
}

// Email は一意なメールアドレスを返す。
func Email() string {
	return New("user") + "@example.com"
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
