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

/// Create response
#[derive(Debug, Serialize)]
pub struct Create {
    /// Sender
    from: H160,
    /// Value
    value: U256,
    /// Gas
    gas: U256,
    /// Initialization code
    init: Bytes,
}

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
	// matched, e := re.FindStringMatch(str)
	// if e != nil {
	// 	t.Fatal(e)
	// }
	matched := FindStruct(str, "Call")
	fmt.Println(matched)
}
