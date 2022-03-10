package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

// func TestParseRusttraitsfile(t *testing.T) {
// 	file := `use cfx_types::{H160, H256, U256, U64};
// use jsonrpc_core::Result;
// use jsonrpc_derive::rpc;

// use crate::rpc::types::{
//     eth::{
//         Block, BlockNumber, CallRequest, EthRpcLogFilter, FilterChanges, Log,
//         Receipt, SyncStatus, Transaction,
//     },
//     Bytes, Index,
// };
// /// Eth rpc interface.
// #[rpc(server)]
// pub trait Eth {
//     #[rpc(name = "web3_clientVersion")]
//     fn client_version(&self) -> Result<String>;
// }
// /// Eth rpc interface.
// #[rpc(server)]
// pub trait Eth {
//     #[rpc(name = "web3_clientVersion")]
//     fn client_version(&self) -> Result<String>;
// }
// /// Eth filters rpc api (polling).
// // TODO: do filters api properly
// #[rpc(server)]
// pub trait EthFilter {
//     /// Returns id of new filter.
//     #[rpc(name = "eth_newFilter")]
//     fn new_filter(&self, _: EthRpcLogFilter) -> Result<U256>;
// }
// #[rpc(server)]
// pub trait EthFilter {
//     /// Returns id of new filter.
//     #[rpc(name = "eth_newFilter")]
//     fn new_filter(&self, _: EthRpcLogFilter) -> Result<U256>;
// }`

// 	traitReg := regexp.MustCompile(`(?mUs)(\/\/\/.*\n|)#\[rpc\(.*\)\]\npub trait .* \{[\s\S]*}`)
// 	traitRegFinded := traitReg.FindAllString(file, -1)
// 	fmt.Printf("traitRegFinded %v\n", traitRegFinded)
// }

func TestParseRusttraitsfile(t *testing.T) {
	var file RustTraitsFile = `
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
	var rt RustTrait = `/// Eth rpc interface.
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
