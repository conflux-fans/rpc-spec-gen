package rust

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseStruct(t *testing.T) {
	var rustStruct Struct = `
	/// Block representation
	// #[derive(Debug, Serialize)]
	// #[serde(rename_all = "camelCase")]
	pub struct Block {
	    /// Hash of the block
	    pub hash: Option<H256>,
	    /// Hash of the parent
	    pub parent_hash: H256,
	}`
	parsed := rustStruct.Parse()

	b, _ := json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s\n", b)

	rustStruct = `
pub struct Account {
	///
	pub address: H256,
	///
	pub block_number: U64,
	///
	pub status: NodeLockStatus,
}`
	parsed = rustStruct.Parse()

	b, _ = json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s\n", b)
}

func TestParseFieldType(t *testing.T) {
	var fieldType RustType = `Option<Vec<U256>>`
	parsed := fieldType.Parse()

	b, _ := json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s\n", b)
}
