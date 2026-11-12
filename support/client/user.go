package client

import (
	"context"
	"net/http"
)

// NOTE: パス・項目名は仮。実際の API 仕様に合わせて修正する。

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserClient struct{ c *httpClient }

func (u *UserClient) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	var out User
	if err := u.c.do(ctx, http.MethodPost, "/users", req, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (u *UserClient) Get(ctx context.Context, id string) (*User, error) {
	var out User
	if err := u.c.do(ctx, http.MethodGet, "/users/"+id, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}
