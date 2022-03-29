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
	structs, us := NewSouceCode(str[1]).GetStructs()
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
	enum, us := NewSouceCode(str).GetEnums()
	for k, v := range enum {
		fmt.Printf("enum k:%v, \nv:%v\n", k, v)
	}
	for k, v := range us {
		fmt.Printf("use k:%v, \nv:%v\n", k, v)
	}
}

func TestNewScurceCode(t *testing.T) {
	raws := []string{`
    // Copyright 2021 Conflux Foundation. All rights reserved.
// Conflux is free software and distributed under GNU General Public License.
// See http://www.gnu.org/licenses/

// Copyright 2015-2020 Parity Technologies (UK) Ltd.
// This file is part of OpenEthereum.

// OpenEthereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// OpenEthereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with OpenEthereum.  If not, see <http://www.gnu.org/licenses/>.

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

    /// Returns true if client is actively mining new blocks.
    #[rpc(name = "eth_mining")]
    fn is_mining(&self) -> Result<bool>;

    /// Returns the chain ID used for transaction signing at the
    /// current best block. None is returned if not
    /// available.
    #[rpc(name = "eth_chainId")]
    fn chain_id(&self) -> Result<Option<U64>>;

    /// Returns current gas_price.
    #[rpc(name = "eth_gasPrice")]
    fn gas_price(&self) -> Result<U256>;

    /// Returns current max_priority_fee
    #[rpc(name = "eth_maxPriorityFeePerGas")]
    fn max_priority_fee_per_gas(&self) -> Result<U256>;

    // /// Returns transaction fee history.
    // #[rpc(name = "eth_feeHistory")]
    // fn fee_history(&self, _: U256, _: BlockNumber, _: Option<Vec<f64>>)
    //     -> BoxFuture<EthFeeHistory>;

    /// Returns balance of the given account.
    #[rpc(name = "eth_getBalance")]
    fn balance(&self, _: H160, _: Option<BlockNumber>) -> Result<U256>;

    // /// Returns the account- and storage-values of the specified account
    // including the Merkle-proof #[rpc(name = "eth_getProof")]
    // fn proof(&self, _: H160, _: Vec<H256>, _: Option<BlockNumber>) ->
    // BoxFuture<EthAccount>;

    /// Returns content of the storage at given address.
    #[rpc(name = "eth_getStorageAt")]
    fn storage_at(
        &self, _: H160, _: U256, _: Option<BlockNumber>,
    ) -> jsonrpc_core::Result<H256>;

    /// Returns an uncles at given block and index.
    #[rpc(name = "eth_getUncleByBlockNumberAndIndex")]
    fn uncle_by_block_number_and_index(
        &self, _: BlockNumber, _: Index,
    ) -> Result<Option<Block>>;

    // /// Returns available compilers.
    // /// @deprecated
    // #[rpc(name = "eth_getCompilers")]
    // fn compilers(&self) -> Result<Vec<String>>;
    //
    // /// Compiles lll code.
    // /// @deprecated
    // #[rpc(name = "eth_compileLLL")]
    // fn compile_lll(&self, _: String) -> Result<Bytes>;
    //
    // /// Compiles solidity.
    // /// @deprecated
    // #[rpc(name = "eth_compileSolidity")]
    // fn compile_solidity(&self, _: String) -> Result<Bytes>;
    //
    // /// Compiles serpent.
    // /// @deprecated
    // #[rpc(name = "eth_compileSerpent")]
    // fn compile_serpent(&self, _: String) -> Result<Bytes>;

    /// Returns logs matching given filter object.
    #[rpc(name = "eth_getLogs")]
    fn logs(&self, _: EthRpcLogFilter) -> Result<Vec<Log>>;

    // /// Returns the hash of the current block, the seedHash, and the boundary
    // condition to be met. #[rpc(name = "eth_getWork")]
    // fn work(&self, _: Option<u64>) -> Result<Work>;

    // /// Used for submitting a proof-of-work solution.
    // #[rpc(name = "eth_submitWork")]
    // fn submit_work(&self, _: H64, _: H256, _: H256) -> Result<bool>;

    /// Used for submitting mining hashrate.
    #[rpc(name = "eth_submitHashrate")]
    fn submit_hashrate(&self, _: U256, _: H256) -> Result<bool>;
}

/// Eth filters rpc api (polling).
// TODO: do filters api properly
#[rpc(server)]
pub trait EthFilter {
    /// Returns id of new filter.
    #[rpc(name = "eth_newFilter")]
    fn new_filter(&self, _: EthRpcLogFilter) -> Result<U256>;

    /// Returns id of new block filter.
    #[rpc(name = "eth_newBlockFilter")]
    fn new_block_filter(&self) -> Result<U256>;

    /// Returns id of new block filter.
    #[rpc(name = "eth_newPendingTransactionFilter")]
    fn new_pending_transaction_filter(&self) -> Result<U256>;

    /// Returns filter changes since last poll.
    #[rpc(name = "eth_getFilterChanges")]
    fn filter_changes(&self, _: Index) -> Result<FilterChanges>;

    /// Returns all logs matching given filter (in a range 'from' - 'to').
    #[rpc(name = "eth_getFilterLogs")]
    fn filter_logs(&self, _: Index) -> Result<Vec<Log>>;

    /// Uninstalls filter.
    #[rpc(name = "eth_uninstallFilter")]
    fn uninstall_filter(&self, _: Index) -> Result<bool>;
}


#[cfg(test)]
mod tests {
    use super::{SyncInfo, SyncStatus};
    use serde_json;

    #[test]
    fn test_serialize_sync_info() {
        let t = SyncInfo::default();
        let serialized = serde_json::to_string(&t).unwrap();
        assert_eq!(
            serialized,
            r#"{"startingBlock":"0x0","currentBlock":"0x0","highestBlock":"0x0","warpChunksAmount":null,"warpChunksProcessed":null}"#
        );
    }
}


/// Transaction
#[derive(Debug, Default, Clone, PartialEq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct Transaction {
    // /// transaction type
    // #[serde(rename = "type", skip_serializing_if = "Option::is_none")]
    // pub transaction_type: Option<U64>,
    /// Hash
    pub hash: H256,
    /// Nonce
    pub nonce: U256,
    /// Block hash
    pub block_hash: Option<H256>,
    /// Block number
    pub block_number: Option<U256>,
    /// Transaction Index
    pub transaction_index: Option<U256>,
    /// Sender
    pub from: H160,
    /// Recipient
    pub to: Option<H160>,
    /// Transfered value
    pub value: U256,
    /// Gas Price
    pub gas_price: U256,
    /// Max fee per gas
    #[serde(skip_serializing_if = "Option::is_none")]
    pub max_fee_per_gas: Option<U256>,
    /// Gas
    pub gas: U256,
    /// Data
    pub input: Bytes,
    /// Creates contract
    pub creates: Option<H160>,
    /// Raw transaction data
    pub raw: Bytes,
    /// Public key of the signer.
    pub public_key: Option<H512>,
    /// The network id of the transaction, if any.
    pub chain_id: Option<U64>,
    /// The standardised V field of the signature (0 or 1). Used by legacy
    /// transaction
    #[serde(skip_serializing_if = "Option::is_none")]
    pub standard_v: Option<U256>,
    /// The standardised V field of the signature.
    pub v: U256,
    /// The R field of the signature.
    pub r: U256,
    /// The S field of the signature.
    pub s: U256,
    // Whether tx is success
    pub status: Option<U64>,
    /* /// Transaction activates at specified block.
     * pub condition: Option<TransactionCondition>,
     * /// optional access list
     * #[serde(skip_serializing_if = "Option::is_none")]
     * pub access_list: Option<AccessList>,
     * /// miner bribe
     * #[serde(skip_serializing_if = "Option::is_none")]
     * pub max_priority_fee_per_gas: Option<U256>, */
}

    `}
	sc := NewSouceCode(raws[0])
	fmt.Println(sc.Cleaned())
}
