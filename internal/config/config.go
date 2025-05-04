package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Schema string
	Output string
}

func (c Config) Validate() error {
	if c.Schema == "" {
		return fmt.Errorf("schema file path is required")
	}

	file, err := os.Open(c.Schema)
	if err != nil {
		return fmt.Errorf("failed to open schema file: %w", err)
	}
	defer file.Close()

	var jsonData map[string]any
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonData); err != nil {
		return fmt.Errorf("invalid json in schema file: %w", err)
	}

	version, ok := jsonData["openapi"].(string)
	if !ok {
		return fmt.Errorf("schema file is not a valid openapi specification - missing 'openapi' version field")
	}

	if !strings.HasPrefix(version, "3.") {
		return fmt.Errorf("schema file must be openapi v3, got version %s", version)
	}

	return nil
}
