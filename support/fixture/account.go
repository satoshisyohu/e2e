package fixture

import (
	"testing"

	"e2e/support/client"
)

// Account はユーザーを作成し、そのユーザーの口座を1件作成して返す。
func (f *Fixture) Account(t testing.TB) (*client.User, *client.Account) {
	t.Helper()

	u := f.User(t)
	acc, err := f.clients.Account.Create(ctx(), client.CreateAccountRequest{UserID: u.ID})
	if err != nil {
		t.Fatalf("fixture: create account: %v", err)
	}
	return u, acc
}
