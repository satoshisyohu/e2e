//go:build e2e

package application_test

import (
	"context"
	"testing"

	"e2e/support/client"
)

func TestCancel(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	_, app := fx.RegisteredApplication(t)

	canceled, err := clients.Application.Cancel(ctx, app.ID)
	if err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if canceled.Status != client.ApplicationStatusCanceled {
		t.Errorf("status = %q, want %q", canceled.Status, client.ApplicationStatusCanceled)
	}
}
