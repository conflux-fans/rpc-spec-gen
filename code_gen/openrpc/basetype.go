package openrpc

import "github.com/go-openapi/spec"

var basetypeSchemas = map[string]spec.Schema{
	"bool": {
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"boolean"},
		},
	},
	"H160": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{40}$",
		},
	},
	"H256": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{64}$",
		},
	},
	"U256": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x([1-9a-f][0-9a-f]{0,63}|0)$",
		},
	},
	"U64": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x([1-9a-f][0-9a-f]{0,15}|0)$",
		},
	},
}
