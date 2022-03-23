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

	rustStruct = `
	#[derive(PartialEq, Debug, Serialize, Deserialize, Eq, Hash, Clone)]
#[serde(rename_all = "camelCase", deny_unknown_fields)]
pub struct CfxRpcLogFilter {
    /// Search will be applied from this epoch number.
    pub from_epoch: Option<EpochNumber>,

    /// Till this epoch number.
    pub to_epoch: Option<EpochNumber>,

    /// Search will be applied from this block number.
    pub from_block: Option<U64>,

    /// Till this block number.
    pub to_block: Option<U64>,

    /// Search will be applied in these blocks if given.
    /// This will override from/to_epoch fields.
    pub block_hashes: Option<Vec<H256>>,

    /// Search addresses.
    ///
    /// If None, match all.
    /// If specified, log must be produced by one of these addresses.
    pub address: Option<VariadicValue<RpcAddress>>,

    /// Search topics.
    ///
    /// Logs can have 4 topics: the function signature and up to 3 indexed
    /// event arguments. The elements of topics match the corresponding
    /// log topics. Example: ["0xA", null, ["0xB", "0xC"], null] matches
    /// logs with "0xA" as the 1st topic AND ("0xB" OR "0xC") as the 3rd
    /// topic. If None, match all.
    pub topics: Option<Vec<VariadicValue<H256>>>,

    /// Logs offset
    ///
    /// If None, return all logs
    /// If specified, should skip the *last* n logs.
    pub offset: Option<U64>,

    /// Logs limit
    ///
    /// If None, return all logs
    /// If specified, should only return *last* n logs
    /// after the offset has been applied.
    pub limit: Option<U64>,
}`
	parsed = rustStruct.Parse()

	b, _ = json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s\n", b)

	rustStruct = `
#[derive(Debug, Serialize, Clone, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct RewardInfo {
    block_hash: H256,
    author: RpcAddress,
    total_reward: U256,
    base_reward: U256,
    tx_fee: U256,
}
	`
	parsed = rustStruct.Parse()

	b, _ = json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s\n", b)
}

func TestParseFieldType(t *testing.T) {
	fieldType := []RustType{
		`U256`,
		`Option<Vec<U256>>`,
		`Vec<Option<U256>>`,
		`BoxFuture<Vec<RpcLog>>`,
	}

	parsed := fieldType[3].Parse()

	b, _ := json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s\n", b)
}

func TestJsonMarshal(t *testing.T) {
	j, _ := json.Marshal("<>")
	fmt.Printf("%s\n", j)
}
