package openrpc

import (
	"time"

	"github.com/Conflux-Chain/rpc-gen/parser/rust"
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
	"Bytes": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: `^0x[0-9a-f]*$`,
		},
	},
	"StorageRoot": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{64}$",
		},
	},
	"PosBlockId": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{64}$",
		},
	},
	"Bloom": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{1024}$",
		},
	},
	"Address": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: `^(NET\d+|CFX|CFXTEST)(:TYPE\..*|):[ABCDEFGHJKMNPRSTUVWXYZ0123456789]42)$`,
		},
	},
	"Account": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: `^(NET\d+|CFX|CFXTEST)(:TYPE\..*|):[ABCDEFGHJKMNPRSTUVWXYZ0123456789]42)$`,
		},
	},
	"EpochNumber": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^.*$",
		},
	},

	// FIXME: 1. fix pattern, 2. basetypeschem 和 customSchemas 中只需要留一个
	"VariadicValue%3CH256%3E": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^.*$",
		},
	},
	"VariadicValue%3CRpcAddress%3E": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^.*$",
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
	// "BlockNumber": {
	// 	SchemaProps: spec.SchemaProps{
	// 		OneOf: []spec.Schema{
	// 			{
	// 				SchemaProps: spec.SchemaProps{
	// 					Title: "Block number",
	// 					Ref:   spec.MustCreateRef(schemaRefRoot + "U64"),
	// 				},
	// 			},
	// 			{
	// 				SchemaProps: spec.SchemaProps{
	// 					Title: "Block tag",
	// 					Enum:  []interface{}{"earliest", "latest_committed", "latest_voted"},
	// 				},
	// 			},
	// 		},
	// 	},
	// },
	// "EpochNumber": {},

	// // FIXME: 1. fix pattern, 2. basetypeschem 和 customSchemas 中只需要留一个
	// "VariadicValue%3CRpcAddress%3E": {
	// 	SchemaProps: spec.SchemaProps{
	// 		Type:    spec.StringOrArray{"string"},
	// 		Pattern: "^.*$",
	// 	},
	// },
	// "VariadicValue%3CH256%3E": {
	// 	SchemaProps: spec.SchemaProps{
	// 		Type:    spec.StringOrArray{"string"},
	// 		Pattern: "^.*$",
	// 	},
	// },
}

func mustGetBasetypeSchemasByUseType(useType rust.UseType) *spec.Schema {

	if useType.Name == "RpcAddress" {
		time.Sleep(0)
	}

	meta, ok := rust.GetUseTypeMeta(useType)
	if !ok {
		logger.Panicf("meta is nil for useType: %s", useType.String())
	}
	if meta.IsBaseType() {
		if v, ok := basetypeSchemas[useType.Name]; ok {
			return v
		}
		logger.Panicf("basetype schemas not found for useType: %s", useType.String())
	}
	logger.Panicf("useType is not basetype: %s, alias %v", useType.String(), useType.Alias)
	return nil
}
