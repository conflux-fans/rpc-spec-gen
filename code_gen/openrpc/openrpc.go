package openrpc

import "github.com/Conflux-Chain/rpc-gen/parser"

type Schema struct {
	REF        *string
	Title      *string
	Pattern    *string
	Type       *string
	Properties map[string]Schema
	Items      *Schema
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

type Doc struct {
	OpenRpc    string
	Info       Info
	Methods    []Method
	Components Components
}

func GenSchema(rustStruct parser.RustStructParsed) Schema {
	panic("not implemented")
}

func GenMethod(rustFunc parser.RustFuncParsed) Method {
	panic("not implemented")
}

func GenMethods(rustTrait parser.RustTraitParsed) []Method {
	panic("not implemented")
}

func GenDoc(rustTraits []parser.RustTraitParsed) []Doc {
	panic("not implemented")
}
