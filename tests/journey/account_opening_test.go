//go:build e2e

package journey_test

import (
	"context"
	"fmt"
	"testing"

	"e2e/support/assertion"
	"e2e/support/client"
)

// ユーザー作成 → 申込 → 承認 → 口座開設 → 入金 までの一連の流れを確認する。
func TestJourney_AccountOpening(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	u := fx.User(t)

	app, err := clients.Application.Register(ctx, client.RegisterApplicationRequest{UserID: u.ID})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if _, err := clients.Application.Approve(ctx, app.ID); err != nil {
		t.Fatalf("approve: %v", err)
	}

	// 承認後、口座は非同期で作成される。
	var acc client.Account
	assertion.Eventually(t, func() error {
		list, err := clients.Account.ListByUser(ctx, u.ID)
		if err != nil {
			return err
		}
		if len(list.Items) != 1 {
			return fmt.Errorf("accounts = %d, want 1", len(list.Items))
		}
		acc = list.Items[0]
		return nil
	}, assertion.WithTimeout(cfg.Eventually.Timeout), assertion.WithInterval(cfg.Eventually.Interval))

	if _, err := clients.Account.Deposit(ctx, acc.ID, client.DepositRequest{Amount: 5000}); err != nil {
		t.Fatalf("deposit: %v", err)
	}

	got, err := clients.Account.Get(ctx, acc.ID)
	if err != nil {
		t.Fatalf("get account: %v", err)
	}
	if got.Balance != 5000 {
		t.Errorf("balance = %d, want 5000", got.Balance)
	}
}
