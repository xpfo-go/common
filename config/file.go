package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadFile(path string, out any) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return LoadJSON(path, out)
	case ".yaml", ".yml":
		return LoadYAML(path, out)
	default:
		return fmt.Errorf("unsupported config file extension: %s", ext)
	}
}

func LoadJSON(path string, out any) error {
	bs, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read json config: %w", err)
	}
	if err := json.Unmarshal(bs, out); err != nil {
		return fmt.Errorf("unmarshal json config: %w", err)
	}
	return nil
}

func LoadYAML(path string, out any) error {
	bs, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read yaml config: %w", err)
	}
	if err := yaml.Unmarshal(bs, out); err != nil {
		return fmt.Errorf("unmarshal yaml config: %w", err)
	}
	return nil
}
