//go:build e2e

package account_test

import (
	"context"
	"testing"

	"e2e/support/client"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	u := fx.User(t)

	acc, err := clients.Account.Create(ctx, client.CreateAccountRequest{UserID: u.ID})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if acc.ID == "" {
		t.Fatal("account ID is empty")
	}
	if acc.UserID != u.ID {
		t.Errorf("user_id = %q, want %q", acc.UserID, u.ID)
	}
	if acc.Balance != 0 {
		t.Errorf("balance = %d, want 0", acc.Balance)
	}
}
