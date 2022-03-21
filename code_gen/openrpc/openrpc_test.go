package openrpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Conflux-Chain/rpc-gen/parser/rust"
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
