package openrpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-openapi/spec"
)

func TestGetRelatedSchemas(t *testing.T) {
	// s := mustLoadSchema("cfx_space", rust.UseType{
	// 	ModPath: []string{"crate", "rpc", "types", "pos"},
	// 	Name:    "Account",
	// })

	content := `
	{
		"type": "object",
		"properties": {
		  "address": {
			"title": "address",
			"$ref": "#/components/schemas/__H256"
		  },
		  "block_number": {
			"title": "block_number",
			"$ref": "#/components/schemas/__U64"
		  },
		  "status": {
			"title": "status",
			"$ref": "#/components/schemas/crate__rpc__types__pos__NodeLockStatus"
		  }
		}
	  }
	`

	s := spec.Schema{}
	if e := json.Unmarshal([]byte(content), &s); e != nil {
		panic(e)
	}

	related := getRelatedSchemas(s, "cfx_space")
	j, _ := json.MarshalIndent(related, "", "  ")
	fmt.Printf("related %s\n", j)
}
