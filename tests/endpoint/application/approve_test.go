//go:build e2e

package application_test

import (
	"context"
	"fmt"
	"testing"

	"e2e/support/assertion"
	"e2e/support/client"
)

// 申込を承認すると、口座サービスに口座が作成されることを確認する。
func TestApprove_CreatesAccount(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	u, app := fx.RegisteredApplication(t)

	approved, err := clients.Application.Approve(ctx, app.ID)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if approved.Status != client.ApplicationStatusApproved {
		t.Errorf("status = %q, want %q", approved.Status, client.ApplicationStatusApproved)
	}

	assertion.Eventually(t, func() error {
		list, err := clients.Account.ListByUser(ctx, u.ID)
		if err != nil {
			return err
		}
		if len(list.Items) != 1 {
			return fmt.Errorf("accounts = %d, want 1", len(list.Items))
		}
		return nil
	}, assertion.WithTimeout(cfg.Eventually.Timeout), assertion.WithInterval(cfg.Eventually.Interval))
}
