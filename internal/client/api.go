package client

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/geokaralis/apigen/internal/config"
)

func Generate(ctx context.Context, config config.Config) error {
	const apitmpl = `/**
 * This source file is auto-generated. Do not edit directly.
 */`

	if config.Output != "" {
		if err := os.MkdirAll(config.Output, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	tmpl, err := template.New("api").Parse(apitmpl)

	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(config.Output + "/api.ts")
	if err != nil {
		return fmt.Errorf("failed to create api.ts file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, apitmpl); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
