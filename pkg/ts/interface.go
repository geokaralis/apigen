package ts

import (
	"fmt"
	"strings"

	"slices"

	"github.com/geokaralis/apigen/pkg/openapi/v3"
	"github.com/geokaralis/apigen/pkg/utils"
)

type InterfaceCompiler struct {
	*TsCompiler
}

func NewInterfaceCompiler(schema *openapi.OpenApiSchema) *InterfaceCompiler {
	return &InterfaceCompiler{
		TsCompiler: New(schema),
	}
}

func (c *InterfaceCompiler) Generate() (string, error) {
	var result strings.Builder
	processedParams := make(map[string]bool)
	paths := c.Schema.Paths
	schemas := c.Schema.Components.Schemas

	for _, pathItem := range paths {
		for _, operation := range pathItem {
			if len(operation.Parameters) > 0 {
				var queryParams []openapi.Parameter
				for _, param := range operation.Parameters {
					if param.In == "query" {
						queryParams = append(queryParams, param)
					}
				}

				if len(queryParams) > 0 {
					interfaceName := fmt.Sprintf("%sParams", utils.ToPascalCase(operation.OperationID))

					if _, exists := processedParams[interfaceName]; !exists {
						processedParams[interfaceName] = true

						if operation.Description != "" {
							result.WriteString(fmt.Sprintf("/**\n * Parameters for %s\n * %s\n */\n",
								operation.OperationID, operation.Description))
						}

						result.WriteString(fmt.Sprintf("export interface %s {\n", interfaceName))

						for _, param := range queryParams {
							fieldType := TypeFromSchema(param.Schema, schemas)

							if param.Description != "" {
								result.WriteString(fmt.Sprintf("  /** %s */\n", param.Description))
							}

							if param.Required {
								result.WriteString(fmt.Sprintf("  %s: %s;\n", param.Name, fieldType))
							} else {
								result.WriteString(fmt.Sprintf("  %s?: %s;\n", param.Name, fieldType))
							}
						}

						result.WriteString("}\n")
					}
				}
			}

			for statusCode, response := range operation.Responses {
				// (2xx)
				if strings.HasPrefix(statusCode, "2") && len(response.Content) > 0 {
					if content, ok := response.Content["application/json"]; ok {
						responseType := ResponseType(content.Schema, schemas)

						if !strings.HasPrefix(responseType, "#") && !strings.Contains(responseType, "[]") {
							interfaceName := fmt.Sprintf("%sResponse", utils.ToPascalCase(operation.OperationID))

							if _, exists := processedParams[interfaceName]; !exists {
								processedParams[interfaceName] = true

								if content.Schema.Ref == "" && content.Schema.Type == "object" && len(content.Schema.Properties) > 0 {
									if response.Description != "" {
										result.WriteString(fmt.Sprintf("/**\n * %s\n */\n", response.Description))
									}

									result.WriteString(fmt.Sprintf("export interface %s {\n", interfaceName))

									for propName, propSchema := range content.Schema.Properties {
										fieldType := TypeFromSchema(propSchema, schemas)
										required := slices.Contains(content.Schema.Required, propName)

										if propSchema.Description != "" {
											result.WriteString(fmt.Sprintf("  /** %s */\n", propSchema.Description))
										}

										if required {
											result.WriteString(fmt.Sprintf("  %s: %s;\n", propName, fieldType))
										} else {
											result.WriteString(fmt.Sprintf("  %s?: %s;\n", propName, fieldType))
										}
									}

									result.WriteString("}\n")
								}
							}
						}
					}
				}
			}
		}
	}

	return result.String(), nil
}
