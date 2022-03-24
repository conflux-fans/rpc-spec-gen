package openrpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Conflux-Chain/jsonrpc-spec/tools/spec-gen/parser/rust"
	"github.com/Conflux-Chain/jsonrpc-spec/tools/spec-gen/utils"
	"github.com/go-openapi/spec"
)

func TestGenSchemas(t *testing.T) {

	us := []rust.UseType{
		{
			ModPath: []string{"crate", "rpc", "types", "pos"},
			Name:    "Block",
		},
	}

	schemas := GenSchemas(us)
	j, _ := json.MarshalIndent(schemas, "", "  ")
	fmt.Printf("schemas %s\n", j)
}

func TestFindUseType(t *testing.T) {
	us := []rust.UseType{
		{
			ModPath: []string{"cfx_types"},
			Name:    "H256",
		},
	}
	mustFindUseType("H256", us)
}

func TestGenObjRefSchema(t *testing.T) {
	_type := rust.TypeParsed{
		IsOption: true,
		Core: &rust.TypeParsed{
			IsArray: true,
			Core: &rust.TypeParsed{
				IsOption: true,
				Core: &rust.TypeParsed{
					Name: "U64",
				},
			},
		},
	}

	s := genObjRefSchema(_type, nil)
	j, _ := json.MarshalIndent(s, "", "  ")
	fmt.Printf("s %s\n", j)
}

func TestGenSchemaForParsedType(t *testing.T) {
	refSchema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef("#/components/schemas/U64"),
		},
	}

	rustTypeParsed := rust.RustType("Option<Vec<U64>>").Parse()
	rustTypeParsed = rust.RustType("Vec<Option<U64>>").Parse()

	s := genSchemaForParsedType(rustTypeParsed, refSchema)
	fmt.Printf("s %+v\n", utils.MustJsonPretty(s))
}

func TestGenSchemaByEnum(t *testing.T) {
	var e rust.Enum = `#[derive(Debug, PartialEq, Clone, Hash, Eq)]
	pub enum BlockHashOrEpochNumber {
		BlockHash(H256),
		EpochNumber(EpochNumber),
	}`

	r := e.Parse()
	s := GenSchemaByEnum(r, nil)
	fmt.Printf("s %+v\n", utils.MustJsonPretty(s))
}
