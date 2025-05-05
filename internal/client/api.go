package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/geokaralis/apigen/internal/config"
	"github.com/geokaralis/apigen/pkg/openapi/v3"
	"github.com/geokaralis/apigen/pkg/ts"
)

func Generate(ctx context.Context, cfg config.Config) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	openapi := openapi.New()
	schema, err := openapi.Parse(cfg.Schema)
	if err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	formatter := ts.NewFormatter(2, true, true)

	tc := ts.NewTypeCompiler(schema)
	types, err := tc.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate typescript types: %w", err)
	}

	ic := ts.NewInterfaceCompiler(schema)
	interfaces, err := ic.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate TypeScript interfaces: %w", err)
	}

	cc := ts.NewClientCompiler(schema)
	methods, err := cc.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate TypeScript client methods: %w", err)
	}

	const ftmpl = `/**
* This source file is auto-generated.
*/

{{.Types}}

{{.Interfaces}}

export class ApiClient {
	private baseUrl: string;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

{{.Methods}}
}`

	tmpl, err := template.New("api").Parse(ftmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buffer strings.Builder
	data := struct {
		Types      string
		Interfaces string
		Methods    string
	}{
		Types:      types,
		Interfaces: interfaces,
		Methods:    methods,
	}

	if err := tmpl.Execute(&buffer, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	ft := formatter.Format(buffer.String())

	out := filepath.Join(cfg.Output, "api.ts")
	if err := os.WriteFile(out, []byte(ft), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
