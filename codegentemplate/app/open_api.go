package app

import "grest.dev/grest"

func OpenAPI() *openAPIUtil {
	if openAPI == nil {
		openAPI = &openAPIUtil{}
	}
	return openAPI
}

var openAPI *openAPIUtil

type openAPIUtil struct {
	grest.OpenAPI
}

func (o *openAPIUtil) Configure() *openAPIUtil {
	o.SetVersion()
	o.Servers = []map[string]any{
		{"description": "Local", "url": "http://localhost:4001"},
	}
	o.Info.Title = "My App API"
	o.Info.Description = ""
	o.Info.Version = APP_VERSION
	if o.Components == nil {
		o.Components = map[string]any{}
	}
	if o.Components["parameters"] == nil {
		o.Components["parameters"] = map[string]any{}
	}
	param, _ := o.Components["parameters"].(map[string]any)
	param["pathParam.ID"] = map[string]any{
		"in":          "path",
		"name":        "id",
		"description": "An ID of the resources",
		"schema":      map[string]any{"type": "string"},
		"required":    true,
	}
	param["queryParam.Any"] = map[string]any{
		"in":   "query",
		"name": "params",
		"schema": map[string]any{
			"type": "object",
			"additionalProperties": map[string]any{
				"type": "string",
			},
		},
		"explode": true,
	}
	param["headerParam.Accept-Language"] = map[string]any{
		"in":   "header",
		"name": "Accept-Language",
		"schema": map[string]any{
			"type":    "string",
			"default": "en-US",
			"enum":    []string{"en-US", "en", "id-ID", "id"},
		},
	}
	o.Components["parameters"] = param
	o.Components["securitySchemes"] = map[string]any{
		"bearerTokenAuth": map[string]any{
			"type":   "http",
			"scheme": "bearer",
		},
	}
	o.Security = []map[string]any{
		{"bearerTokenAuth": []string{}},
	}
	return o
}

type OpenAPIOperationInterface interface {
	grest.OpenAPIOperationInterface
}

type OpenAPIOperation struct {
	grest.OpenAPIOperation
}

func OpenAPIError() *openAPIError {
	return &openAPIError{}
}

type openAPIError struct {
	StatusCode  int
	Message     string
	SchemaName  string
	Description string
	Headers     map[string]any
	Links       map[string]any
}

func (o *openAPIError) BadRequest() map[string]any {
	o.StatusCode = 400
	o.Message = "The request cannot be performed because of malformed or missing parameters."
	o.SchemaName = "Error.BadRequest"
	o.Description = "A validation exception has occurred."
	return o.Response()
}

func (o *openAPIError) Unauthorized() map[string]any {
	o.StatusCode = 401
	o.Message = "Invalid authentication token."
	o.SchemaName = "Error.Unauthorized"
	o.Description = "Invalid authorization credentials."
	return o.Response()
}

func (o *openAPIError) Forbidden() map[string]any {
	o.StatusCode = 403
	o.Message = "The user does not have permission to access the resource."
	o.SchemaName = "Error.Forbidden"
	o.Description = "User doesn't have permission to access the resource."
	return o.Response()
}

func (o *openAPIError) Response() map[string]any {
	res := map[string]any{
		"content": map[string]any{
			"application/json": o,
		},
	}
	if o.Description != "" {
		res["description"] = o.Description
	}
	if len(o.Headers) > 0 {
		res["headers"] = o.Headers
	}
	if len(o.Links) > 0 {
		res["links"] = o.Links
	}
	return res
}

func (o *openAPIError) OpenAPISchemaName() string {
	return o.SchemaName
}

func (o *openAPIError) GetOpenAPISchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"error": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"code": map[string]any{
						"type":    "integer",
						"format":  "int32",
						"example": o.StatusCode,
					},
					"message": map[string]any{
						"type":    "string",
						"example": o.Message,
					},
				},
			},
		},
	}
}
