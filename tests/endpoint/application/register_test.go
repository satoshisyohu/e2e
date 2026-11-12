//go:build e2e

package application_test

import (
	"context"
	"testing"

	"e2e/support/client"
)

func TestRegister(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	u := fx.User(t)

	app, err := clients.Application.Register(ctx, client.RegisterApplicationRequest{UserID: u.ID})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if app.ID == "" {
		t.Fatal("application ID is empty")
	}
	if app.Status != client.ApplicationStatusRegistered {
		t.Errorf("status = %q, want %q", app.Status, client.ApplicationStatusRegistered)
	}
}
