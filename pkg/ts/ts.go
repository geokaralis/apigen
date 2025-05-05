package ts

import (
	"fmt"
	"strings"

	"slices"

	"github.com/geokaralis/apigen/pkg/openapi/v3"
)

type Compiler interface {
	Generate() (string, error)
}

type TsCompiler struct {
	Schema *openapi.OpenApiSchema
}

func New(schema *openapi.OpenApiSchema) *TsCompiler {
	return &TsCompiler{
		Schema: schema,
	}
}

func Type(schema openapi.SchemaRef, schemas map[string]openapi.SchemaObject) string {
	if schema.Ref != "" {
		parts := strings.Split(schema.Ref, "/")
		return parts[len(parts)-1]
	}

	switch schema.Type {
	case "string":
		if len(schema.Enum) > 0 {
			var unionParts []string
			for _, enumValue := range schema.Enum {
				unionParts = append(unionParts, fmt.Sprintf("'%s'", enumValue))
			}
			return strings.Join(unionParts, " | ")
		}
		return "string"
	case "integer", "number":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		if schema.Items != nil {
			itemType := Type(*schema.Items, schemas)
			return itemType + "[]"
		}
		return "any[]"
	case "object":
		if schema.Properties != nil && len(schema.Properties) > 0 {
			var objType strings.Builder
			objType.WriteString("{ ")

			for propName, propSchema := range schema.Properties {
				fieldType := Type(propSchema, schemas)
				required := slices.Contains(schema.Required, propName)

				if required {
					objType.WriteString(fmt.Sprintf("%s: %s; ", propName, fieldType))
				} else {
					objType.WriteString(fmt.Sprintf("%s?: %s; ", propName, fieldType))
				}
			}

			objType.WriteString("}")
			return objType.String()
		}
		return "Record<string, any>"
	default:
		return "any"
	}
}

func TypeFromSchema(schema openapi.SchemaRef, schemas map[string]openapi.SchemaObject) string {
	return Type(schema, schemas)
}

func ResponseType(schema openapi.SchemaRef, schemas map[string]openapi.SchemaObject) string {
	return Type(schema, schemas)
}
