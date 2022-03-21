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

func mustGetBasetypeSchemasByUseType(useType rust.UseType) *spec.Schema {
	meta := config.GetUseTypeMeta(useType)
	if meta == nil {
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
