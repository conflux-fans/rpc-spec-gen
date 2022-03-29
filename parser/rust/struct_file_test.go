package rust

import (
	"fmt"
	"testing"
)

func TestCutStructFromFile(t *testing.T) {
	// var re = regexp2.MustCompile(`\/\/\/(?:.(?!\/\/\/))+pub struct Call {.*?}`, regexp2.Multiline|regexp2.Singleline)
	// // re2 := regexp2.MustCompile("^error((?!timeout).)*$", regexp2.RE2)
	var str = `use crate::rpc::types::{
    Action as RpcCfxAction, Bytes, LocalizedTrace as RpcCfxLocalizedTrace,
};
use cfx_types::{H160, H256, U256};
use cfxcore::{observer::trace::Outcome, vm::CallType as CfxCallType};
use jsonrpc_core::Error as JsonRpcError;
use serde::{ser::SerializeStruct, Serialize, Serializer};
use std::{
    convert::{TryFrom, TryInto},
    fmt,
};

/// Call type.
#[derive(Debug, Serialize)]
#[serde(rename_all = "lowercase")]
pub enum CallType {
    /// None
    None,
    /// Call
    Call,
    /// Call code
    CallCode,
    /// Delegate call
    DelegateCall,
    /// Static call
    StaticCall,
}

impl From<CfxCallType> for CallType {
    fn from(cfx_call_type: CfxCallType) -> Self {
        match cfx_call_type {
            CfxCallType::None => CallType::None,
            CfxCallType::Call => CallType::Call,
            CfxCallType::CallCode => CallType::CallCode,
            CfxCallType::DelegateCall => CallType::DelegateCall,
            CfxCallType::StaticCall => CallType::StaticCall,
        }
    }
}

/// Call response
#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Call {
    /// Sender
    from: H160,
    /// Recipient
    to: H160,
    /// Transfered Value
    value: U256,
    /// Gas
    gas: U256,
    /// Input data
    input: Bytes,
    /// The type of the call.
    call_type: CallType,
}
}`
	matched, uses := NewSouceCode(str).FindStruct("Call")
	fmt.Println(matched)
	fmt.Println(uses)

	str = `// Copyright 2019 Conflux Foundation. All rights reserved.
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
    `
	matched, uses = NewSouceCode(str).FindStruct("Account")
	fmt.Println(matched)
	fmt.Println(uses)

	str = `
// Copyright 2019 Conflux Foundation. All rights reserved.
// Conflux is free software and distributed under GNU General Public License.
// See http://www.gnu.org/licenses/

use super::Decision;
use cfx_types::{H256, U64};
use serde::Serialize;

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct Block {
    ///
    pub hash: H256,
    ///
    pub height: U64,
    ///
    pub epoch: U64,
    ///
    pub round: U64,
    ///
    pub last_tx_number: U64,
    ///
    pub miner: Option<H256>,
    ///
    pub parent_hash: H256,
    ///
    pub timestamp: U64,
    ///
    pub pivot_decision: Option<Decision>,
    ///
    pub signatures: Vec<Signature>,
}

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct Signature {
    ///
    pub account: H256,
    ///
    // pub signature: String,
    ///
    pub votes: U64,
}
    `
	matched, uses = NewSouceCode(str).FindStruct("Block")
	fmt.Println(matched)
	fmt.Println(uses)
}
