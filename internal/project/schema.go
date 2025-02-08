package project

// JSONSchemaType represents a JSON Schema type
type JSONSchemaType string

const (
	JSONSchemaTypeString  JSONSchemaType = "string"
	JSONSchemaTypeNumber  JSONSchemaType = "number"
	JSONSchemaTypeInteger JSONSchemaType = "integer"
	JSONSchemaTypeObject  JSONSchemaType = "object"
	JSONSchemaTypeArray   JSONSchemaType = "array"
	JSONSchemaTypeBoolean JSONSchemaType = "boolean"
	JSONSchemaTypeNull    JSONSchemaType = "null"
)

// JSONSchemaDefinition represents a JSON Schema definition
type JSONSchemaDefinition struct {
	Type        JSONSchemaType                  `json:"type,omitempty"`
	Description string                          `json:"description,omitempty"`
	Properties  map[string]JSONSchemaDefinition `json:"properties,omitempty"`
	Items       *JSONSchemaDefinition           `json:"items,omitempty"`
	Required    []string                        `json:"required,omitempty"`
	Enum        []string                        `json:"enum,omitempty"`
	Ref         string                          `json:"$ref,omitempty"`
	OneOf       []JSONSchemaDefinition          `json:"oneOf,omitempty"`
	AnyOf       []JSONSchemaDefinition          `json:"anyOf,omitempty"`
	AllOf       []JSONSchemaDefinition          `json:"allOf,omitempty"`
	Not         *JSONSchemaDefinition           `json:"not,omitempty"`
	Definitions map[string]JSONSchemaDefinition `json:"definitions,omitempty"`
}
