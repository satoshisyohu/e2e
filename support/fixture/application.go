package fixture

import (
	"testing"

	"e2e/support/client"
)

// RegisteredApplication はユーザーを作成し、登録済みの申込を1件作成して返す。
func (f *Fixture) RegisteredApplication(t testing.TB) (*client.User, *client.Application) {
	t.Helper()

	u := f.User(t)
	app, err := f.clients.Application.Register(ctx(), client.RegisterApplicationRequest{UserID: u.ID})
	if err != nil {
		t.Fatalf("fixture: register application: %v", err)
	}
	return u, app
}
