// Package assertion は E2E テスト向けのアサーションを提供する。
package assertion

import (
	"testing"
	"time"
)

type options struct {
	timeout  time.Duration
	interval time.Duration
}

// Option は Eventually の挙動を変更する。
type Option func(*options)

// WithTimeout は待機の上限時間を指定する。
func WithTimeout(d time.Duration) Option { return func(o *options) { o.timeout = d } }

// WithInterval は再試行の間隔を指定する。
func WithInterval(d time.Duration) Option { return func(o *options) { o.interval = d } }

// Eventually は check が nil を返すまで、interval ごとに再試行する。
// timeout を過ぎても成功しなければ、最後のエラーを添えてテストを失敗させる。
//
// 非同期処理(キューやイベント経由のサービス間連携)の結果を待つときに使う。
// 固定の time.Sleep は使わないこと。
func Eventually(t testing.TB, check func() error, opts ...Option) {
	t.Helper()

	o := options{timeout: 30 * time.Second, interval: 1 * time.Second}
	for _, opt := range opts {
		opt(&o)
	}

	deadline := time.Now().Add(o.timeout)
	var lastErr error
	for attempt := 1; ; attempt++ {
		if lastErr = check(); lastErr == nil {
			return
		}
		if time.Now().Add(o.interval).After(deadline) {
			t.Fatalf("condition not met within %s (%d attempts): %v", o.timeout, attempt, lastErr)
		}
		time.Sleep(o.interval)
	}
}
