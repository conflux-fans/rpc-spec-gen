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
