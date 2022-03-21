package rust

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dlclark/regexp2"
)

func TestParseTraitsFile(t *testing.T) {
	content := `
use crate::rpc::types::pos::{
	Account, Block, BlockNumber, CommitteeState, PoSEpochReward, Status,
	Transaction,
};
use cfx_types::{H256, U64};
use diem_types::{
	epoch_state::EpochState, ledger_info::LedgerInfoWithSignatures,
};
use jsonrpc_core::Result as JsonRpcResult;
use jsonrpc_derive::rpc;

/// PoS specific rpc interface.
#[rpc(server)]
pub trait Pos {
	#[rpc(name = "pos_getStatus")]
    fn pos_status(&self) -> JsonRpcResult<Status>;

    #[rpc(name = "pos_getAccount")]
    fn pos_account(
        &self, address: H256, view: Option<U64>,
    ) -> JsonRpcResult<Account>;
}`
	p := TraitsFile(content).Parse()

	j, _ := json.MarshalIndent(p, "", " ")

	fmt.Printf("%v\n", string(j))
}

func TestRegex2(t *testing.T) {
	var re = regexp2.MustCompile(`\/\/\/(?:.(?!\/\/\/))+pub struct ([^\{]*) \{.*?}|pub struct ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	str := `// Copyright 2019 Conflux Foundation. All rights reserved.
// Conflux is free software and distributed under GNU General Public License.
// See http://www.gnu.org/licenses/

use crate::rpc::types::pos::NodeLockStatus;
use cfx_types::{H256, U64};
use serde_derive::Serialize;

#[derive(Debug, Serialize, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct Account {
    ///
    pub address: H256,
    ///
    pub block_number: U64,
    ///
    pub status: NodeLockStatus,
}


/// i am account
#[derive(Debug, Serialize, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct Account {
    ///
    pub address: H256,
    ///
    pub block_number: U64,
    ///
    pub status: NodeLockStatus,
}
	`
	m, _ := re.FindStringMatch(str)
	fmt.Printf("capture groups len %v, %v\n,%v\n,%v\n\n", len(m.Groups()),
		m.Groups()[0].Capture.String(),
		m.Groups()[1].Capture.String(),
		m.Groups()[2].Capture.String(),
	)
	// m, _ = re.FindNextMatch(m)
	// fmt.Printf("capture %v\n", m.Groups()[1].Capture.String())
}

func TestGetStructs(t *testing.T) {
	str := `// Copyright 2019 Conflux Foundation. All rights reserved.
// Conflux is free software and distributed under GNU General Public License.
// See http://www.gnu.org/licenses/

use crate::rpc::types::pos::NodeLockStatus;
use cfx_types::{H256, U64};
use serde_derive::Serialize;

#[derive(Debug, Serialize, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct Account {
    ///
    pub address: H256,
    ///
    pub block_number: U64,
    ///
    pub status: NodeLockStatus,
}


/// i am account
#[derive(Debug, Serialize, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct Account {
    ///
    pub address: H256,
    ///
    pub block_number: U64,
    ///
    pub status: NodeLockStatus,
}
	`
	structs, us := GetStructs(str)
	for k, v := range structs {
		fmt.Printf("struct k:%v, \nv:%v\n", k, v)
	}
	for k, v := range us {
		fmt.Printf("use k:%v, \nv:%v\n", k, v)
	}
}

func TestGetEnums(t *testing.T) {
	str := `// Copyright 2019 Conflux Foundation. All rights reserved.
// Conflux is free software and distributed under GNU General Public License.
// See http://www.gnu.org/licenses/

use super::Decision;
use cfx_types::{H256, U64};
use serde_derive::Serialize;

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Status {
	///
	pub latest_committed: U64,
	///
	pub epoch: U64,
	///
	pub pivot_decision: Decision,
	///
	pub latest_voted: Option<U64>,
	///
	pub latest_tx_number: U64,
}

impl Default for Status {
	fn default() -> Status {
		Status {
			epoch: U64::default(),
			latest_committed: U64::default(),
			pivot_decision: Decision {
				height: U64::default(),
				block_hash: H256::default(),
			},
			latest_voted: None,
			latest_tx_number: U64::default(),
		}
	}
}`
	enum, us := GetEnums(str)
	for k, v := range enum {
		fmt.Printf("enum k:%v, \nv:%v\n", k, v)
	}
	for k, v := range us {
		fmt.Printf("use k:%v, \nv:%v\n", k, v)
	}
}
