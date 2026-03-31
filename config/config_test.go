package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testConfig struct {
	AppName string        `json:"app_name" yaml:"app_name" default:"demo-app"`
	Port    int           `json:"port" yaml:"port" default:"8080"`
	Debug   bool          `json:"debug" yaml:"debug"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" default:"3s"`
	Tags    []string      `json:"tags" yaml:"tags" default:"api,core"`
	DB      struct {
		Host string `json:"host" yaml:"host" required:"true"`
		User string `json:"user" yaml:"user" default:"root"`
	} `json:"db" yaml:"db"`
}

func TestLoadYAMLApplyDefaultsEnvAndValidate(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yaml")
	content := []byte("app_name: file-app\nport: 9000\ndb:\n  host: file-db\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	var cfg testConfig
	if err := LoadFile(path, &cfg); err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if err := ApplyDefaults(&cfg); err != nil {
		t.Fatalf("ApplyDefaults() error = %v", err)
	}

	t.Setenv("APP_PORT", "9100")
	t.Setenv("APP_DB_HOST", "env-db")
	t.Setenv("APP_DEBUG", "true")
	if err := LoadFromEnv("APP", &cfg); err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}
	if err := ValidateRequired(&cfg); err != nil {
		t.Fatalf("ValidateRequired() error = %v", err)
	}

	if cfg.AppName != "file-app" {
		t.Fatalf("AppName = %q, want file-app", cfg.AppName)
	}
	if cfg.Port != 9100 {
		t.Fatalf("Port = %d, want 9100", cfg.Port)
	}
	if cfg.DB.Host != "env-db" {
		t.Fatalf("DB.Host = %q, want env-db", cfg.DB.Host)
	}
	if cfg.DB.User != "root" {
		t.Fatalf("DB.User = %q, want root", cfg.DB.User)
	}
	if cfg.Timeout != 3*time.Second {
		t.Fatalf("Timeout = %s, want 3s", cfg.Timeout)
	}
	if len(cfg.Tags) != 2 || cfg.Tags[0] != "api" || cfg.Tags[1] != "core" {
		t.Fatalf("Tags = %#v, want [api core]", cfg.Tags)
	}
	if !cfg.Debug {
		t.Fatalf("Debug = false, want true")
	}
}

func TestValidateRequiredFails(t *testing.T) {
	var cfg testConfig
	if err := ApplyDefaults(&cfg); err != nil {
		t.Fatal(err)
	}
	if err := ValidateRequired(&cfg); err == nil {
		t.Fatalf("expected required validation error")
	}
}

func TestLoadJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.json")
	content := []byte(`{"app_name":"json-app","db":{"host":"json-db"}}`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	var cfg testConfig
	if err := LoadJSON(path, &cfg); err != nil {
		t.Fatalf("LoadJSON() error = %v", err)
	}
	if cfg.AppName != "json-app" {
		t.Fatalf("AppName = %q, want json-app", cfg.AppName)
	}
	if cfg.DB.Host != "json-db" {
		t.Fatalf("DB.Host = %q, want json-db", cfg.DB.Host)
	}
}
