package client

import (
	"context"
	"net/http"
)

// NOTE: パス・項目名・ステータス値は仮。実際の API 仕様に合わせて修正する。

const (
	ApplicationStatusRegistered = "registered"
	ApplicationStatusApproved   = "approved"
	ApplicationStatusCanceled   = "canceled"
)

type RegisterApplicationRequest struct {
	UserID string `json:"user_id"`
}

type Application struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type ApplicationClient struct{ c *httpClient }

func (a *ApplicationClient) Register(ctx context.Context, req RegisterApplicationRequest) (*Application, error) {
	var out Application
	if err := a.c.do(ctx, http.MethodPost, "/applications", req, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *ApplicationClient) Approve(ctx context.Context, id string) (*Application, error) {
	var out Application
	if err := a.c.do(ctx, http.MethodPost, "/applications/"+id+"/approve", nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *ApplicationClient) Cancel(ctx context.Context, id string) (*Application, error) {
	var out Application
	if err := a.c.do(ctx, http.MethodPost, "/applications/"+id+"/cancel", nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *ApplicationClient) Get(ctx context.Context, id string) (*Application, error) {
	var out Application
	if err := a.c.do(ctx, http.MethodGet, "/applications/"+id, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}
