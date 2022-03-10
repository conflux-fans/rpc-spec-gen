package openrpc

import "github.com/Conflux-Chain/rpc-gen/parser/rust"

type Doc struct {
	OpenRpc    string
	Info       Info
	Methods    []Method
	Components Components
}
type Schema struct {
	REF        *string           `json:"$ref,omitempty"`
	Title      *string           `json:"title,omitempty"`
	Pattern    *string           `json:"pattern,omitempty"`
	Type       *string           `json:"type,omitempty"`
	Required   []string          `json:"required,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
	OneOf      []Schema          `json:"oneOf,omitempty"`
}

type Param struct {
	Name     string
	Required bool
	// maybe OpenRpcSchema{REF} or OpenRpcSchema{others}
	Schema Schema
}

type Result struct {
	Name string
	// maybe OpenRpcSchema{REF} or OpenRpcSchema{others}
	Schema Schema
}

type Info struct {
	Title       string
	Description string
	Version     string
	License     License
}

type License struct {
	Name string
	URL  string
}

type Method struct {
	Name    string
	Summary string
	Params  []Param
	Result  Result
}

type Components struct {
	Schemas map[string]Schema
}

func GenSchema(structParsed rust.StructParsed) Schema {
	panic("not implemented")
}

func GenMethod(funcParsed rust.FuncParsed) Method {
	panic("not implemented")
}

func GenMethods(traitParsed rust.TraitParsed) []Method {
	panic("not implemented")
}

func GenDoc(fileParsed []rust.TraitsFileParsed) []Doc {
	panic("not implemented")
}
