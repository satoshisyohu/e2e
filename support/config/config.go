// Package config は E2E テストの実行設定を envs/<env>.yaml から読み込む。
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Service は接続先サービスごとの設定。
type Service struct {
	BaseURL string `yaml:"base_url"`
}

// Config は E2E テスト全体の設定。
type Config struct {
	Env string `yaml:"-"`

	Services struct {
		User        Service `yaml:"user"`
		Application Service `yaml:"application"`
		Account     Service `yaml:"account"`
	} `yaml:"services"`

	HTTPTimeout time.Duration `yaml:"http_timeout"`

	Eventually struct {
		Timeout  time.Duration `yaml:"timeout"`
		Interval time.Duration `yaml:"interval"`
	} `yaml:"eventually"`
}

// Load は環境変数 E2E_ENV(既定値: local)に対応する envs/<env>.yaml を読み込む。
// サービスの接続先は環境変数 E2E_<SERVICE>_BASE_URL で上書きできる。
func Load() (*Config, error) {
	env := os.Getenv("E2E_ENV")
	if env == "" {
		env = "local"
	}

	path := filepath.Join(envsDir(), env+".yaml")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	cfg := &Config{Env: env}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	override(&cfg.Services.User, "USER")
	override(&cfg.Services.Application, "APPLICATION")
	override(&cfg.Services.Account, "ACCOUNT")
	setDefaults(cfg)

	return cfg, nil
}

func override(s *Service, name string) {
	if v := os.Getenv("E2E_" + strings.ToUpper(name) + "_BASE_URL"); v != "" {
		s.BaseURL = v
	}
}

func setDefaults(cfg *Config) {
	if cfg.HTTPTimeout == 0 {
		cfg.HTTPTimeout = 10 * time.Second
	}
	if cfg.Eventually.Timeout == 0 {
		cfg.Eventually.Timeout = 30 * time.Second
	}
	if cfg.Eventually.Interval == 0 {
		cfg.Eventually.Interval = 1 * time.Second
	}
}

// envsDir はテストの実行ディレクトリに依存せず envs/ を解決する。
// E2E_CONFIG_DIR が指定されていればそれを優先する。
func envsDir() string {
	if d := os.Getenv("E2E_CONFIG_DIR"); d != "" {
		return d
	}
	_, file, _, _ := runtime.Caller(0)
	// support/config/config.go -> リポジトリルート
	return filepath.Join(filepath.Dir(file), "..", "..", "envs")
}
