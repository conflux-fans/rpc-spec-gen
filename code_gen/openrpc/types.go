package openrpc

import (
	"github.com/go-openapi/spec"
)

type OpenRPCSpec1 struct {
	OpenRPC      string        `json:"openrpc,omitempty"`
	Info         Info          `json:"info,omitempty"`
	Servers      []*Server     `json:"servers,omitempty"`
	Methods      []*Method     `json:"methods,omitempty"`
	Components   *Components   `json:"components,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type Info struct {
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Server struct {
	Name        string                    `json:"name,omitempty"`
	URL         string                    `json:"url,omitempty"`
	Summary     string                    `json:"summary,omitempty"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Method struct {
	Name           string               `json:"name,omitempty"`
	Tags           []Tag                `json:"tags,omitempty"`
	Summary        string               `json:"summary,omitempty"`
	Description    string               `json:"description,omitempty"`
	ExternalDocs   *ExternalDocs        `json:"externalDocs,omitempty"`
	Params         []*ContentDescriptor `json:"params,omitempty"`
	Result         *ContentDescriptor   `json:"result,omitempty"`
	Deprecated     bool                 `json:"deprecated,omitempty"`
	Servers        []Server             `json:"servers,omitempty"`
	Errors         []Error              `json:"errors,omitempty"`
	Links          []Link               `json:"links,omitempty"`
	ParamStructure string               `json:"paramStructure,omitempty"`
	Examples       []ExamplePairing     `json:"examples,omitempty"`
}

type ContentDescriptor struct {
	Content
}
type Content struct {
	Name        string      `json:"name,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Deprecated  bool        `json:"deprecated,omitempty"`
	Schema      spec.Schema `json:"schema,omitempty"`
}

type Example struct {
	Name          string      `json:"name,omitempty"`
	Summary       string      `json:"summary,omitempty"`
	Description   string      `json:"description,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty"`
}

type ExamplePairing struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	Params      []Example `json:"params,omitempty"`
	Result      Example   `json:"result,omitempty"`
}

type Link struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Method      string                 `json:"method,omitempty"`
	Params      map[string]interface{} `json:"params,omitempty"`
	Server      Server                 `json:"server,omitempty"`
}

// https://www.jsonrpc.org/specification#error_object
type Error struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Components struct {
	ContentDescriptors    map[string]*ContentDescriptor `json:"contentDescriptors,omitempty"`
	Schemas               map[string]*spec.Schema       `json:"schemas,omitempty"`
	Examples              map[string]Example            `json:"examples,omitempty"`
	Links                 map[string]Link               `json:"links,omitempty"`
	Errors                map[string]Error              `json:"errors,omitempty"`
	ExamplePairingObjects map[string]ExamplePairing     `json:"examplePairingObjects,omitempty"`
	Tags                  map[string]Tag                `json:"tags,omitempty"`
}

type Tag struct {
	Name         string        `json:"name,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// type Schema struct {
// 	Ref        *string           `json:"$ref,omitempty"`
// 	Title      *string           `json:"title,omitempty"`
// 	Pattern    *string           `json:"pattern,omitempty"`
// 	Type       *string           `json:"type,omitempty"`
// 	Required   []string          `json:"required,omitempty"`
// 	Properties map[string]Schema `json:"properties,omitempty"`
// 	Items      *Schema           `json:"items,omitempty"`
// 	OneOf      []Schema          `json:"oneOf,omitempty"`
// }

// type Doc struct {
// 	OpenRpc    string
// 	Info       Info
// 	Methods    []Method
// 	Components Components
// }

// type Param struct {
// 	Name     string
// 	Required bool
// 	// maybe OpenRpcSchema{REF} or OpenRpcSchema{others}
// 	Schema Schema
// }

// type Result struct {
// 	Name string
// 	// maybe OpenRpcSchema{REF} or OpenRpcSchema{others}
// 	Schema Schema
// }

// type Info struct {
// 	Title       string
// 	Description string
// 	Version     string
// 	License     License
// }

// type License struct {
// 	Name string
// 	URL  string
// }

// type Method struct {
// 	Name    string
// 	Summary string
// 	Params  []Param
// 	Result  Result
// }

// type Components struct {
// 	Schemas map[string]Schema
// }
