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

	rustStruct = `
/// The validator node returns this structure which includes signatures
/// from validators that confirm the state.  The client needs to only pass back
/// the LedgerInfo element since the validator node doesn't need to know the
/// signatures again when the client performs a query, those are only there for
/// the client to be able to verify the state
#[derive(Clone, Debug, Eq, PartialEq, Serialize, Deserialize)]
pub struct LedgerInfoWithV0 {
	ledger_info: LedgerInfo,
	/// The validator is identified by its account address: in order to verify
	/// a signature one needs to retrieve the public key of the validator
	/// for the given epoch.
	signatures: BTreeMap<AccountAddress, ConsensusSignature>,
}
	`
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
