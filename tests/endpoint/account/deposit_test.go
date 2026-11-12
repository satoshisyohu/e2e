//go:build e2e

package account_test

import (
	"context"
	"testing"

	"e2e/support/client"
)

func TestDeposit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	_, acc := fx.Account(t)

	if _, err := clients.Account.Deposit(ctx, acc.ID, client.DepositRequest{Amount: 1000}); err != nil {
		t.Fatalf("deposit: %v", err)
	}

	got, err := clients.Account.Get(ctx, acc.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Balance != 1000 {
		t.Errorf("balance = %d, want 1000", got.Balance)
	}
}
