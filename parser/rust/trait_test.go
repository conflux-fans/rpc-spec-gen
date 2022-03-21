package rust

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseRusttraitsfile(t *testing.T) {
	var file TraitsFile = `
//! Eth rpc interface.
use cfx_types::{H160, H256, U256, U64};
use jsonrpc_core::Result;
use jsonrpc_derive::rpc;

use crate::rpc::types::{
	eth::{
		Block, BlockNumber, CallRequest, EthRpcLogFilter, FilterChanges, Log,
		Receipt, SyncStatus, Transaction,
	},
	Bytes, Index,
};

/// Eth rpc interface.
#[rpc(server)]
pub trait Eth {
	#[rpc(name = "web3_clientVersion")]
	fn client_version(&self) -> Result<String>;

	#[rpc(name = "net_version")]
	fn net_version(&self) -> Result<String>;

	/// Returns protocol version encoded as a string (quotes are necessary).
	#[rpc(name = "eth_protocolVersion")]
	fn protocol_version(&self) -> Result<String>;
}`

	parsed := file.Parse()
	j, _ := json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s", j)
}

func TestParseRustTrait(t *testing.T) {
	var rt Trait = `/// Eth rpc interface.
#[rpc(server)]
pub trait Eth {
	#[rpc(name = "web3_clientVersion")]
	fn client_version(&self) -> Result<String>;

	#[rpc(name = "net_version")]
	fn net_version(&self) -> Result<String>;

	/// Returns protocol version encoded as a string (quotes are necessary).
	#[rpc(name = "eth_protocolVersion")]
	fn protocol_version(&self) -> Result<String>;
}`
	parsed := rt.Parse()
	j, _ := json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s", j)

	rt = `
#[rpc(server)]
pub trait Eth {
	#[rpc(name = "web3_clientVersion")]
	fn client_version(&self) -> Result<String>;

	#[rpc(name = "net_version")]
	fn net_version(&self) -> Result<String>;

	/// Returns protocol version encoded as a string (quotes are necessary).
	#[rpc(name = "eth_protocolVersion")]
	fn protocol_version(&self) -> Result<String>;
}`
	parsed = rt.Parse()
	j, _ = json.MarshalIndent(parsed, "", "  ")
	fmt.Printf("%s", j)

}
