package openapi

type OpenApiSchema struct {
	OpenApi    string              `json:"openapi"`
	Paths      map[string]PathItem `json:"paths"`
	Components struct {
		Schemas map[string]SchemaObject `json:"schemas"`
	} `json:"components"`
}

type PathItem map[string]OperationObject

type OperationObject struct {
	OperationID string              `json:"operationId"`
	Parameters  []Parameter         `json:"parameters"`
	Responses   map[string]Response `json:"responses"`
	Summary     string              `json:"summary"`
	Description string              `json:"description"`
}

type Parameter struct {
	Name        string    `json:"name"`
	In          string    `json:"in"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
	Schema      SchemaRef `json:"schema"`
}

type Response struct {
	Description string                   `json:"description"`
	Content     map[string]ContentObject `json:"content"`
}

type ContentObject struct {
	Schema SchemaRef `json:"schema"`
}

type SchemaRef struct {
	Ref         string               `json:"$ref"`
	Type        string               `json:"type"`
	Items       *SchemaRef           `json:"items"`
	Properties  map[string]SchemaRef `json:"properties"`
	Required    []string             `json:"required"`
	Description string               `json:"description"`
	Format      string               `json:"format"`
	Enum        []string             `json:"enum"`
	Default     any                  `json:"default"`
}

type SchemaObject struct {
	Type        string               `json:"type"`
	Required    []string             `json:"required"`
	Properties  map[string]SchemaRef `json:"properties"`
	Enum        []string             `json:"enum"`
	OneOf       []SchemaRef          `json:"oneOf"`
	Description string               `json:"description"`
	Format      string               `json:"format"`
}
