package fixture

import (
	"testing"

	"e2e/support/client"
	"e2e/support/idgen"
)

// User はユーザーを1件作成して返す。
func (f *Fixture) User(t testing.TB) *client.User {
	t.Helper()

	u, err := f.clients.User.Create(ctx(), client.CreateUserRequest{
		Name:  idgen.New("name"),
		Email: idgen.Email(),
	})
	if err != nil {
		t.Fatalf("fixture: create user: %v", err)
	}
	return u
}
