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
	str := []string{
		`// Copyright 2019 Conflux Foundation. All rights reserved.
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
}`, `
		// Copyright 2019 Conflux Foundation. All rights reserved.
// Conflux is free software and distributed under GNU General Public License.
// See http://www.gnu.org/licenses/

use cfx_types::{H256, U64};
use diem_types::{
    transaction::{ConflictSignature, TransactionPayload, TransactionStatus},
    vm_status::KeptVMStatus,
};
use rustc_hex::ToHex;
use serde::{ser::SerializeStruct, Serialize, Serializer};

#[derive(Debug, Clone)]
pub struct Transaction {
    pub hash: H256,
    pub from: H256,
    pub block_hash: Option<H256>,
    pub block_number: Option<U64>,
    pub timestamp: Option<U64>,
    pub number: U64,
    pub payload: Option<RpcTransactionPayload>,
    pub status: Option<RpcTransactionStatus>,
    pub tx_type: RpcTransactionType,
}

#[derive(Debug, Serialize, Clone, Copy)]
pub enum RpcTransactionStatus {
    Executed,
    Failed,
    Discard,
}

#[derive(Debug, Serialize, Clone, Copy)]
pub enum RpcTransactionType {
    BlockMetadata,
    Election,
    Retire,
    Register,
    UpdateVotingPower,
    PivotDecision,
    Dispute,
    Other,
}

#[derive(Debug, Clone)]
pub enum RpcTransactionPayload {
    ///
    Register(RegisterPayload),
    ///
    Election(ElectionPayload),
    ///
    UpdateVotingPower(UpdateVotingPowerPayload),
    ///
    PivotDecision(PivotDecisionPayload),
    ///
    Retire(RetirePayload),
    ///
    Dispute(DisputePayload),
    ///
    Other,
}

impl From<TransactionPayload> for RpcTransactionPayload {
    fn from(payload: TransactionPayload) -> Self {
        match payload {
            TransactionPayload::Election(e) => {
                RpcTransactionPayload::Election(ElectionPayload {
                    public_key: format!("0x{}", e.public_key),
                    target_term: U64::from(e.target_term),
                    vrf_proof: format!("0x{}", e.vrf_proof),
                    vrf_public_key: format!("0x{}", e.vrf_public_key),
                })
            }
            TransactionPayload::Retire(r) => {
                RpcTransactionPayload::Retire(RetirePayload {
                    address: H256::from(r.node_id.to_u8()),
                    voting_power: U64::from(r.votes),
                })
            }
            TransactionPayload::Register(r) => {
                RpcTransactionPayload::Register(RegisterPayload {
                    vrf_public_key: format!("0x{}", r.vrf_public_key),
                    public_key: format!("0x{}", r.public_key),
                })
            }
            TransactionPayload::UpdateVotingPower(u) => {
                RpcTransactionPayload::UpdateVotingPower(
                    UpdateVotingPowerPayload {
                        address: H256::from(u.node_address.to_u8()),
                        voting_power: U64::from(u.voting_power),
                    },
                )
            }
            TransactionPayload::PivotDecision(p) => {
                RpcTransactionPayload::PivotDecision(PivotDecisionPayload {
                    height: U64::from(p.height),
                    block_hash: H256::from(p.block_hash),
                })
            }
            TransactionPayload::Dispute(d) => {
                let conflicting_votes = match d.conflicting_votes {
                    ConflictSignature::Proposal((first, second)) => {
                        ConflictingVotes {
                            conflict_vote_type: "proposal".into(),
                            first: format!("0x{}", first.to_hex::<String>()),
                            second: format!("0x{}", second.to_hex::<String>()),
                        }
                    }
                    ConflictSignature::Vote((first, second)) => {
                        ConflictingVotes {
                            conflict_vote_type: "vote".into(),
                            first: format!("0x{}", first.to_hex::<String>()),
                            second: format!("0x{}", second.to_hex::<String>()),
                        }
                    }
                };
                RpcTransactionPayload::Dispute(DisputePayload {
                    address: H256::from(d.address.to_u8()),
                    bls_public_key: format!("0x{}", d.bls_pub_key),
                    vrf_public_key: format!("0x{}", d.vrf_pub_key),
                    conflicting_votes,
                })
            }
            _ => RpcTransactionPayload::Other,
        }
    }
}

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct RegisterPayload {
    pub vrf_public_key: String,
    pub public_key: String,
}

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct ElectionPayload {
    pub public_key: String,
    pub target_term: U64,
    pub vrf_proof: String,
    pub vrf_public_key: String,
}

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct UpdateVotingPowerPayload {
    pub address: H256,
    pub voting_power: U64,
}

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct PivotDecisionPayload {
    pub height: U64,
    pub block_hash: H256,
}

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct RetirePayload {
    pub address: H256,
    pub voting_power: U64,
}

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct DisputePayload {
    pub address: H256,
    pub bls_public_key: String,
    pub vrf_public_key: String,
    pub conflicting_votes: ConflictingVotes,
}

#[derive(Debug, Serialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct ConflictingVotes {
    pub conflict_vote_type: String,
    pub first: String,
    pub second: String,
}

impl Serialize for Transaction {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where S: Serializer {
        let mut struc = serializer.serialize_struct("Transaction", 9)?;
        struc.serialize_field("hash", &self.hash)?;
        struc.serialize_field("from", &self.from)?;
        struc.serialize_field("number", &self.number)?;
        struc.serialize_field("blockHash", &self.block_hash)?;
        struc.serialize_field("blockNumber", &self.block_number)?;
        struc.serialize_field("timestamp", &self.timestamp)?;
        struc.serialize_field("status", &self.status)?;
        struc.serialize_field("type", &self.tx_type)?;
        if self.payload.is_some() {
            match &self.payload.as_ref().unwrap() {
                RpcTransactionPayload::Election(e) => {
                    struc.serialize_field("payload", e)?;
                }
                RpcTransactionPayload::Retire(r) => {
                    struc.serialize_field("payload", r)?;
                }
                RpcTransactionPayload::Register(r) => {
                    struc.serialize_field("payload", r)?;
                }
                RpcTransactionPayload::UpdateVotingPower(u) => {
                    struc.serialize_field("payload", u)?;
                }
                RpcTransactionPayload::PivotDecision(p) => {
                    struc.serialize_field("payload", p)?;
                }
                RpcTransactionPayload::Dispute(d) => {
                    struc.serialize_field("payload", d)?;
                }
                _ => {}
            }
        } else {
            let empty: Option<TransactionPayload> = None;
            struc.serialize_field("payload", &empty)?
        }
        struc.end()
    }
}

pub fn tx_type(payload: TransactionPayload) -> RpcTransactionType {
    match payload {
        TransactionPayload::Election(_) => RpcTransactionType::Election,
        TransactionPayload::Retire(_) => RpcTransactionType::Retire,
        TransactionPayload::Register(_) => RpcTransactionType::Register,
        TransactionPayload::UpdateVotingPower(_) => {
            RpcTransactionType::UpdateVotingPower
        }
        TransactionPayload::PivotDecision(_) => {
            RpcTransactionType::PivotDecision
        }
        TransactionPayload::Dispute(_) => RpcTransactionType::Dispute,
        _ => RpcTransactionType::Other,
    }
}

impl From<TransactionStatus> for RpcTransactionStatus {
    fn from(status: TransactionStatus) -> Self {
        match status {
            TransactionStatus::Discard(_) => RpcTransactionStatus::Discard,
            TransactionStatus::Keep(keep_status) => match keep_status {
                KeptVMStatus::Executed => RpcTransactionStatus::Executed,
                _ => RpcTransactionStatus::Failed,
            },
            TransactionStatus::Retry => RpcTransactionStatus::Failed,
        }
    }
}

impl From<KeptVMStatus> for RpcTransactionStatus {
    fn from(status: KeptVMStatus) -> Self {
        match status {
            KeptVMStatus::Executed => RpcTransactionStatus::Executed,
            _ => RpcTransactionStatus::Failed,
        }
    }
}

		`,
	}
	structs, us := SourceCode(str[1]).GetStructs()
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
	enum, us := SourceCode(str).GetEnums()
	for k, v := range enum {
		fmt.Printf("enum k:%v, \nv:%v\n", k, v)
	}
	for k, v := range us {
		fmt.Printf("use k:%v, \nv:%v\n", k, v)
	}
}
