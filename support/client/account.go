package client

import (
	"context"
	"net/http"
	"net/url"
)

// NOTE: パス・項目名は仮。実際の API 仕様に合わせて修正する。

type CreateAccountRequest struct {
	UserID string `json:"user_id"`
}

type DepositRequest struct {
	Amount int64 `json:"amount"`
}

type Account struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Balance int64  `json:"balance"`
}

type AccountList struct {
	Items []Account `json:"items"`
}

type AccountClient struct{ c *httpClient }

func (a *AccountClient) Create(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	var out Account
	if err := a.c.do(ctx, http.MethodPost, "/accounts", req, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *AccountClient) Deposit(ctx context.Context, id string, req DepositRequest) (*Account, error) {
	var out Account
	if err := a.c.do(ctx, http.MethodPost, "/accounts/"+id+"/deposits", req, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *AccountClient) Get(ctx context.Context, id string) (*Account, error) {
	var out Account
	if err := a.c.do(ctx, http.MethodGet, "/accounts/"+id, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *AccountClient) ListByUser(ctx context.Context, userID string) (*AccountList, error) {
	var out AccountList
	path := "/accounts?user_id=" + url.QueryEscape(userID)
	if err := a.c.do(ctx, http.MethodGet, path, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}
