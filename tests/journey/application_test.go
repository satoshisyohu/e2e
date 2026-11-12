//go:build e2e

package journey_test

import (
	"context"
	"testing"

	"e2e/support/client"
)

// 申込を登録して取り消すと、口座が作成されず、申込が取消状態のまま残ることを確認する。
func TestJourney_ApplicationRegisterAndCancel(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	u := fx.User(t)

	app, err := clients.Application.Register(ctx, client.RegisterApplicationRequest{UserID: u.ID})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if _, err := clients.Application.Cancel(ctx, app.ID); err != nil {
		t.Fatalf("cancel: %v", err)
	}

	got, err := clients.Application.Get(ctx, app.ID)
	if err != nil {
		t.Fatalf("get application: %v", err)
	}
	if got.Status != client.ApplicationStatusCanceled {
		t.Errorf("status = %q, want %q", got.Status, client.ApplicationStatusCanceled)
	}

	list, err := clients.Account.ListByUser(ctx, u.ID)
	if err != nil {
		t.Fatalf("list accounts: %v", err)
	}
	if len(list.Items) != 0 {
		t.Errorf("accounts = %d, want 0", len(list.Items))
	}
}
