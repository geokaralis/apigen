package ts

import (
	"fmt"
	"slices"
	"strings"

	"github.com/geokaralis/apigen/pkg/openapi/v3"
)

type TypeCompiler struct {
	*TsCompiler
}

func NewTypeCompiler(schema *openapi.OpenApiSchema) *TypeCompiler {
	return &TypeCompiler{
		TsCompiler: New(schema),
	}
}

func (c *TypeCompiler) Generate() (string, error) {
	var result strings.Builder
	schemas := c.Schema.Components.Schemas

	for name, schema := range schemas {
		if schema.Description != "" {
			result.WriteString(fmt.Sprintf("/**\n * %s\n */\n", schema.Description))
		}

		if schema.Type == "object" {
			result.WriteString(fmt.Sprintf("export type %s = {\n", name))

			for propName, prop := range schema.Properties {
				fieldType := Type(prop, schemas)
				optional := true

				if slices.Contains(schema.Required, propName) {
					optional = false
				}

				if prop.Description != "" {
					result.WriteString(fmt.Sprintf("  /** %s */\n", prop.Description))
				}

				if optional {
					result.WriteString(fmt.Sprintf("  %s?: %s;\n", propName, fieldType))
				} else {
					result.WriteString(fmt.Sprintf("  %s: %s;\n", propName, fieldType))
				}
			}

			result.WriteString("}\n\n")
		} else if schema.Type == "string" && len(schema.Enum) > 0 {
			result.WriteString(fmt.Sprintf("export type %s = %s;\n\n",
				name, Type(openapi.SchemaRef{Type: "string", Enum: schema.Enum}, schemas)))
		}
	}

	return result.String(), nil
}
