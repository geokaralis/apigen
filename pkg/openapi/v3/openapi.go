package openapi

import (
	"encoding/json"
	"fmt"
	"os"
)

type SchemaParser interface {
	Parse(filePath string) (*OpenApiSchema, error)
}

type OpenApiParser struct{}

func New() *OpenApiParser {
	return &OpenApiParser{}
}

func (p *OpenApiParser) Parse(filePath string) (*OpenApiSchema, error) {
	schemaData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	var schema OpenApiSchema
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	return &schema, nil
}
