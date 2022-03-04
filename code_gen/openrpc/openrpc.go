package openrpc

import "github.com/Conflux-Chain/rpc-check/parser"

type OpenRpcSchema struct {
	REF        *string
	Title      *string
	Pattern    *string
	Type       *string
	Properties map[string]OpenRpcSchema
	Items      *OpenRpcSchema
}

type OpenRpcParam struct {
	Name     string
	Required bool
	// maybe OpenRpcSchema{REF} or OpenRpcSchema{others}
	Schema OpenRpcSchema
}

type OpenRpcResult struct {
	Name string
	// maybe OpenRpcSchema{REF} or OpenRpcSchema{others}
	Schema OpenRpcSchema
}

type OpenRpcInfo struct {
	Title       string
	Description string
	Version     string
	License     OpenRpcLicense
}

type OpenRpcLicense struct {
	Name string
	URL  string
}

type OpenRpcMethod struct {
	Name    string
	Summary string
	Params  []OpenRpcParam
	Result  OpenRpcResult
}

type OpenRpcComponents struct {
	Schemas map[string]OpenRpcSchema
}

type OpenRpcDoc struct {
	OpenRpc    string
	Info       OpenRpcInfo
	Methods    []OpenRpcMethod
	Components OpenRpcComponents
}

func GenSchema(rustStruct parser.RustStruct) OpenRpcSchema {
	panic("not implemented")
}

func GenMethod(rustFunc parser.RustFunc) OpenRpcMethod {
	panic("not implemented")
}

func GenMethods(rustTrait parser.RsutTrait) []OpenRpcMethod {
	panic("not implemented")
}

func GenDoc(rustTrait parser.RsutTrait) []OpenRpcDoc {
	panic("not implemented")
}
