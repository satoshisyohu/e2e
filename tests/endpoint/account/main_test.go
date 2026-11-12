//go:build e2e

package account_test

import (
	"fmt"
	"os"
	"testing"

	"e2e/support/client"
	"e2e/support/config"
	"e2e/support/fixture"
)

var (
	cfg     *config.Config
	clients *client.Clients
	fx      *fixture.Fixture
)

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	clients = client.New(cfg)
	fx = fixture.New(clients)
	os.Exit(m.Run())
}
