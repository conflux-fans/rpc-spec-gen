package openrpc

import (
	"github.com/Conflux-Chain/rpc-gen/parser/rust"
	"github.com/Conflux-Chain/rpc-gen/parser/rust/config"
	"github.com/go-openapi/spec"
)

var basetypeSchemas = map[string]*spec.Schema{
	"bool": {
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"boolean"},
		},
	},
	"u64": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: `^[1-9]\d*$`,
		},
	},
	"u8": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: `^[1-9]\d*$`,
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
	"String": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: `^.*$`,
		},
	},
	"RpcAddress": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: `^(NET\d+|CFX|CFXTEST)(:TYPE\..*|):[ABCDEFGHJKMNPRSTUVWXYZ0123456789]42)$`,
		},
	},
}

// "BlockTag": {
// 	"title": "Block tag",
// 	"type": "string",
// 	"enum": ["earliest", "latest", "pending"]
// },
// "BlockNumberOrTag": {
// 	"title": "Block number or tag",
// 	"oneOf": [
// 		{
// 			"title": "Block number",
// 			"$ref": "#/components/schemas/uint"
// 		},
// 		{
// 			"title": "Block tag",
// 			"$ref": "#/components/schemas/BlockTag"
// 		}
// 	]
// }
var customSchemas = map[string]*spec.Schema{
	"BlockNumber": {
		SchemaProps: spec.SchemaProps{
			OneOf: []spec.Schema{
				{
					SchemaProps: spec.SchemaProps{
						Title: "Block number",
						Ref:   spec.MustCreateRef(schemaRefRoot + "U64"),
					},
				},
				{
					SchemaProps: spec.SchemaProps{
						Title: "Block tag",
						Enum:  []interface{}{"earliest", "latest_committed", "latest_voted"},
					},
				},
			},
		},
	},
}

func mustGetBasetypeSchemasByUseType(useType rust.UseType) *spec.Schema {
	meta, ok := config.GetUseTypeMeta(useType)
	if !ok {
		logger.Panicf("meta is nil for useType: %s", useType.String())
	}
	if meta.IsBaseType() {
		if v, ok := basetypeSchemas[useType.Name]; ok {
			return v
		}
		logger.Panicf("basetype schemas not found for useType: %s", useType.String())
	}
	logger.Panicf("useType is not basetype: %s", useType.String())
	return nil
}
