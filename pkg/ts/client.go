package ts

import (
	"fmt"
	"strings"

	"github.com/geokaralis/apigen/pkg/openapi/v3"
	"github.com/geokaralis/apigen/pkg/utils"
)

type ClientCompiler struct {
	*TsCompiler
}

func NewClientCompiler(schema *openapi.OpenApiSchema) *ClientCompiler {
	return &ClientCompiler{
		TsCompiler: New(schema),
	}
}

func (c *ClientCompiler) Generate() (string, error) {
	var result strings.Builder
	paths := c.Schema.Paths

	for path, pathItem := range paths {
		for method, operation := range pathItem {
			if operation.OperationID == "" {
				continue
			}

			methodName := utils.ToCamelCase(operation.OperationID)

			var methodParams []string
			var pathParams []string
			var queryParams []string

			for _, param := range operation.Parameters {
				if param.In == "path" {
					pathParams = append(pathParams, param.Name)
					if param.Required {
						methodParams = append(methodParams, fmt.Sprintf("%s: string", param.Name))
					} else {
						methodParams = append(methodParams, fmt.Sprintf("%s?: string", param.Name))
					}
				} else if param.In == "query" {
					queryParams = append(queryParams, param.Name)
				}
			}

			if len(queryParams) > 0 {
				methodParams = append(methodParams, fmt.Sprintf("params?: %sParams", utils.ToPascalCase(operation.OperationID)))
			}

			var responseType string
			for statusCode, response := range operation.Responses {
				if strings.HasPrefix(statusCode, "2") && len(response.Content) > 0 {
					if content, ok := response.Content["application/json"]; ok {
						if content.Schema.Ref != "" {
							parts := strings.Split(content.Schema.Ref, "/")
							responseType = parts[len(parts)-1]
						} else if content.Schema.Type == "object" {
							responseType = fmt.Sprintf("%sResponse", utils.ToPascalCase(operation.OperationID))
						} else if content.Schema.Type == "array" && content.Schema.Items != nil {
							if content.Schema.Items.Ref != "" {
								parts := strings.Split(content.Schema.Items.Ref, "/")
								itemType := parts[len(parts)-1]
								responseType = fmt.Sprintf("%s[]", itemType)
							} else {
								responseType = "any[]"
							}
						} else {
							responseType = "any"
						}
						break
					}
				}
			}

			if responseType == "" {
				responseType = "void"
			}

			if operation.Description != "" || operation.Summary != "" {
				result.WriteString("  /**\n")
				if operation.Summary != "" {
					result.WriteString(fmt.Sprintf("   * %s\n", operation.Summary))
				}
				if operation.Description != "" {
					result.WriteString(fmt.Sprintf("   * %s\n", operation.Description))
				}
				result.WriteString("   */\n")
			}

			result.WriteString(fmt.Sprintf("  async %s(%s): Promise<%s> {\n",
				methodName, strings.Join(methodParams, ", "), responseType))

			urlPath := path

			for _, paramName := range pathParams {
				urlPath = strings.Replace(urlPath, fmt.Sprintf("{%s}", paramName), fmt.Sprintf("${%s}", paramName), -1)
			}

			if len(queryParams) > 0 {
				result.WriteString("    const queryParams = new URLSearchParams();\n")
				result.WriteString("    if (params) {\n")
				for _, paramName := range queryParams {
					result.WriteString(fmt.Sprintf("      if (params.%s !== undefined) {\n", paramName))
					result.WriteString(fmt.Sprintf("        queryParams.append('%s', String(params.%s));\n", paramName, paramName))
					result.WriteString("      }\n")
				}
				result.WriteString("    }\n\n")

				result.WriteString(fmt.Sprintf("    const url = `${this.baseUrl}%s${queryParams.toString() ? '?' + queryParams.toString() : ''}`;\n", urlPath))
			} else {
				result.WriteString(fmt.Sprintf("    const url = `${this.baseUrl}%s`;\n", urlPath))
			}

			result.WriteString(fmt.Sprintf("    const response = await fetch(url, {\n"))
			result.WriteString(fmt.Sprintf("      method: '%s',\n", strings.ToUpper(method)))
			result.WriteString("      headers: {\n")
			result.WriteString("        'Content-Type': 'application/json',\n")
			result.WriteString("      },\n")
			result.WriteString("    });\n\n")

			result.WriteString("    if (!response.ok) {\n")
			result.WriteString("      const error = await response.json();\n")
			result.WriteString("      throw new Error(error.message || 'Request failed');\n")
			result.WriteString("    }\n\n")

			if responseType != "void" {
				result.WriteString("    return response.json();\n")
			}

			result.WriteString("  }\n")
		}
	}

	return result.String(), nil
}
