package specconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/conflux-fans/rpc-spec-gen/config"
	"github.com/conflux-fans/rpc-spec-gen/parser/rust"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger = &logrus.Logger{
	Out:   os.Stderr,
	Level: logrus.DebugLevel,
	Formatter: &prefixed.TextFormatter{
		// DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	},
}

const (
	SchemaRefRoot = "#/components/schemas/"
)

var BasetypeSchemas = map[string]*spec.Schema{
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
	"f64": {
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"number"},
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
	"H128": {
		SchemaProps: spec.SchemaProps{
			Type:    spec.StringOrArray{"string"},
			Pattern: "^0x[0-9,a-f,A-F]{28}$",
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
			Pattern: `^(NET\d+|CFX|CFXTEST)(:TYPE\..*|):[ABCDEFGHJKMNPRSTUVWXYZ0123456789]{42}$`,
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

var CustomSchemas = map[string]*spec.Schema{
	"crate::rpc::types::eth::BlockNumber": {
		SchemaProps: spec.SchemaProps{
			Title: "Block Number or Tags",
			OneOf: []spec.Schema{
				{
					SchemaProps: spec.SchemaProps{
						Title: "U64",
						Ref:   spec.MustCreateRef(SchemaRefRoot + "U64"),
					},
				},
				{
					SchemaProps: spec.SchemaProps{
						Type:     spec.StringOrArray{"object"},
						Title:    "BlockHash",
						Required: []string{"blockHash"},
						Properties: map[string]spec.Schema{
							"blockHash": {
								SchemaProps: spec.SchemaProps{
									Ref: spec.MustCreateRef(SchemaRefRoot + "H256"),
								},
							},
							"requireCanonical": {
								SchemaProps: spec.SchemaProps{
									Type: spec.StringOrArray{"boolean"},
								},
							},
						},
					},
				},
				{
					SchemaProps: spec.SchemaProps{
						Type:     spec.StringOrArray{"object"},
						Title:    "BlockNumber",
						Required: []string{"blockNumber"},
						Properties: map[string]spec.Schema{
							"blockNumber": {
								SchemaProps: spec.SchemaProps{
									Ref: spec.MustCreateRef(SchemaRefRoot + "U64"),
								},
							},
						},
					},
				},
				{
					SchemaProps: spec.SchemaProps{
						Enum: []interface{}{
							"earliest",
							"latest",
							"pending",
						},
					},
				},
			},
		},
	},
	// "BlockNumber": {
	// 	SchemaProps: spec.SchemaProps{
	// 		OneOf: []spec.Schema{
	// 			{
	// 				SchemaProps: spec.SchemaProps{
	// 					Title: "Block number",
	// 					Ref:   spec.MustCreateRef(SchemaRefRoot + "U64"),
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

var SchemaPostHooks = map[string]func(schema *spec.Schema){
	"crate::rpc::types::eth::Block": func(schema *spec.Schema) {
		var propSchema spec.Schema
		propSchema.Title = "sha3uncle"
		propSchema.Ref = schema.Properties["unclesHash"].Ref
		schema.Properties["sha3uncle"] = propSchema

		delete(schema.Properties, "unclesHash")
	},
}

func MustGetBasetypeSchemasByUseType(useType rust.UseType) *spec.Schema {

	if useType.Name == "RpcAddress" {
		time.Sleep(0)
	}

	meta, ok := rust.GetUseTypeMeta(useType)
	if !ok {
		logrus.Panicf("meta is nil for useType: %s", useType.String())
	}
	if meta.IsBaseType() {
		if v, ok := BasetypeSchemas[useType.Name]; ok {
			return v
		}
		logrus.Panicf("basetype schemas not found for useType: %s", useType.String())
	}
	logrus.Panicf("useType is not basetype: %s, alias %v", useType.String(), useType.Alias)
	return nil
}

func GetSchemaSaveRelativePath(space string, schemaFullName string) string {
	return path.Join(space, strings.Join(strings.Split(schemaFullName, "::"), "/")+".json")
}

func GetSchemaSavePath(space string, schemaFullName string) string {
	return path.Join(config.GetConfig().SchemaRootPath, GetSchemaSaveRelativePath(space, schemaFullName))
}

func GetUseTypeRefSchema(useType rust.UseType) spec.Schema {
	s := spec.Schema{}

	schemaName := strings.Join(useType.ModPath, "__") + "__" + useType.Name
	schemaName = strings.TrimPrefix(schemaName, "__")
	s.Ref = spec.MustCreateRef(SchemaRefRoot + schemaName)
	return s
}

func ParseSchemaRefToUseType(ref string) rust.UseType {

	logrus.WithField("ref", ref).Debug("parse schema ref to use type")

	fullUseType := strings.TrimPrefix(ref, SchemaRefRoot)
	matchs := regexp.MustCompile(`(.*)__(.*)`).FindStringSubmatch(fullUseType)

	if len(matchs) != 3 {
		logrus.WithField("ref", ref).Debug("parse to a base type")
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

func MustLoadSchema(space string, useType rust.UseType) *spec.Schema {

	if useType.IsBaseType() {
		return MustGetBasetypeSchemasByUseType(useType)
	}

	savePath := GetSchemaSavePath(space, useType.String())
	content, e := ioutil.ReadFile(savePath)

	if e != nil {
		panic(e)
	}

	schema := spec.Schema{}
	if e := json.Unmarshal(content, &schema); e != nil {
		panic(e)
	}

	// j, _ := json.MarshalIndent(schema, "", "  ")

	// logrus.WithField("space", space).WithField("useType", useType).WithField("schema", string(j)).Debug("load schema")
	return &schema
}
