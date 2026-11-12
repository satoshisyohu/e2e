// Package client は各サービスの公開 API を呼び出すクライアントを提供する。
//
// NOTE: 現在は手書きの仮実装。各サービスの OpenAPI 定義が揃い次第、
// oapi-codegen などで生成したクライアントに置き換える。
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"e2e/support/config"
	"e2e/support/idgen"
)

// Clients は全サービスのクライアントをまとめたもの。
type Clients struct {
	User        *UserClient
	Application *ApplicationClient
	Account     *AccountClient
}

// New は設定からすべてのクライアントを生成する。
func New(cfg *config.Config) *Clients {
	hc := &http.Client{Timeout: cfg.HTTPTimeout}
	return &Clients{
		User:        &UserClient{c: newHTTP(cfg.Services.User.BaseURL, hc)},
		Application: &ApplicationClient{c: newHTTP(cfg.Services.Application.BaseURL, hc)},
		Account:     &AccountClient{c: newHTTP(cfg.Services.Account.BaseURL, hc)},
	}
}

// StatusError は想定外の HTTP ステータスが返ったことを表す。
// 異常系のテストでは errors.As で取り出してステータスを確認する。
type StatusError struct {
	Method     string
	URL        string
	StatusCode int
	Body       string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("%s %s: unexpected status %d: %s", e.Method, e.URL, e.StatusCode, e.Body)
}

type httpClient struct {
	baseURL string
	hc      *http.Client
}

func newHTTP(baseURL string, hc *http.Client) *httpClient {
	return &httpClient{baseURL: baseURL, hc: hc}
}

// do はリクエストを送り、wantStatus のときに限りレスポンスを out にデコードする。
func (c *httpClient) do(ctx context.Context, method, path string, in, out any, wantStatus int) error {
	url := c.baseURL + path

	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// 各サービスのログから E2E の実行を追跡できるようにする。
	req.Header.Set("X-E2E-Run-Id", idgen.RunID())

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != wantStatus {
		return &StatusError{Method: method, URL: url, StatusCode: resp.StatusCode, Body: string(respBody)}
	}
	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("unmarshal response: %w: %s", err, respBody)
		}
	}
	return nil
}
