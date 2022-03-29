package openrpc

import (
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/conflux-fans/rpc-spec-gen/config"
	"github.com/conflux-fans/rpc-spec-gen/parser/rust"
	"github.com/go-openapi/spec"
)

const (
	schemaRefRoot = "#/components/schemas/"
)

var basetypeSchemas = map[string]*spec.Schema{
	"bool": {
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"boolean"},
		},
	},
	"u64": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"number"},
			Pattern: `^[1-9]\d*$`,
		},
	},
	"u8": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"number"},
			Pattern: `^[1-9]\d*$`,
		},
	},
	"usize": {
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"number"},
		},
	},
	"H64": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{16}$",
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
	"H512": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{128}$",
		},
	},
	"H2048": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{512}$",
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
			Pattern: `^(NET\d+|CFX|CFXTEST)(:TYPE\..*|):[ABCDEFGHJKMNPRSTUVWXYZ0123456789]{42}$`,
		},
	},
	// "EpochNumber": {
	// 	SchemaProps: spec.SchemaProps{
	// 		Type:    spec.StringOrArray{"string"},
	// 		Pattern: "^.*$",
	// 	},
	// },
	"Index": {
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"number"},
		},
	},
}

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

func getSchemaSaveRelativePath(space string, schemaFullName string) string {
	return path.Join(space, strings.Join(strings.Split(schemaFullName, "::"), "/")+".json")
}

func getSchemaSavePath(space string, schemaFullName string) string {
	return path.Join(config.GetConfig().SchemaRootPath, getSchemaSaveRelativePath(space, schemaFullName))
}

func getUseTypeRefSchema(useType rust.UseType) spec.Schema {
	s := spec.Schema{}

	schemaName := strings.Join(useType.ModPath, "__") + "__" + useType.Name
	schemaName = strings.TrimPrefix(schemaName, "__")
	s.Ref = spec.MustCreateRef(schemaRefRoot + schemaName)
	return s
}

func parseSchemaRefToUseType(ref string) rust.UseType {

	fullUseType := strings.TrimPrefix(ref, schemaRefRoot)
	matchs := regexp.MustCompile(`(.*)__(.*)`).FindStringSubmatch(fullUseType)

	if len(matchs) != 3 {
		logger.WithField("ref", ref).Debug("parse to a base type")
		return rust.UseType{
			Name: fullUseType,
		}
	}

	fullName := matchs[1]

	return rust.UseType{
		ModPath: strings.Split(fullName, "__"),
		Name:    matchs[2],
	}
}
