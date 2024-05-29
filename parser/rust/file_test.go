package rust

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/stretchr/testify/assert"
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

func TestGetDefineTypes(t *testing.T) {
	raw := `// Copyright 2019 Conflux Foundation. All rights reserved.
    // Conflux is free software and distributed under GNU General Public License.
    // See http://www.gnu.org/licenses/
    
    pub mod eth_transaction;
    pub mod native_transaction;
    
    pub use eth_transaction::{
        Eip1559Transaction, Eip155Transaction, Eip2930Transaction,
        EthereumTransaction,
    };
    pub use native_transaction::{
        Cip1559Transaction, Cip2930Transaction, NativeTransaction,
        TypedNativeTransaction,
    };
    
    use crate::{bytes::Bytes, hash::keccak};
    use cfx_types::{
        Address, AddressSpaceUtil, AddressWithSpace, BigEndianHash, Space, H160,
        H256, U256,
    };
    use eth_transaction::eip155_signature;
    use keylib::{
        self, public_to_address, recover, verify_public, Public, Secret, Signature,
    };
    use malloc_size_of::{MallocSizeOf, MallocSizeOfOps};
    use rlp::{self, Decodable, DecoderError, Encodable, Rlp, RlpStream};
    use serde::{Deserialize, Serialize};
    use std::{
        error, fmt,
        ops::{Deref, DerefMut},
    };
    use unexpected::OutOfBounds;
    
    /// Fake address for unsigned transactions.
    pub const UNSIGNED_SENDER: Address = H160([0xff; 20]);
    
    pub const TYPED_NATIVE_TX_PREFIX: &[u8; 3] = b"cfx";
    pub const TYPED_NATIVE_TX_PREFIX_BYTE: u8 = TYPED_NATIVE_TX_PREFIX[0];
    pub const EIP2930_TYPE: u8 = 0x01;
    pub const EIP1559_TYPE: u8 = 0x02;
    pub const CIP2930_TYPE: u8 = 0x01;
    pub const CIP1559_TYPE: u8 = 0x02;
    
    /// Shorter id for transactions in compact blocks
    // TODO should be u48
    pub type TxShortId = u64;
    
    pub type TxPropagateId = u32;
    
    // FIXME: Most errors here are bounded for TransactionPool and intended for rpc,
    // FIXME: however these are unused, they are not errors for transaction itself.
    // FIXME: Transaction verification and consensus related error can be separated.
    #[derive(Debug, PartialEq, Clone)]
    /// Errors concerning transaction processing.
    pub enum TransactionError {
        /// Transaction is already imported to the queue
        AlreadyImported,
        /// Chain id in the transaction doesn't match the chain id of the network.
        ChainIdMismatch {
            expected: u32,
            got: u32,
            space: Space,
        },
        /// Epoch height out of bound.
        EpochHeightOutOfBound {
            block_height: u64,
            set: u64,
            transaction_epoch_bound: u64,
        },
        /// The gas paid for transaction is lower than base gas.
        NotEnoughBaseGas {
            /// Absolute minimum gas required.
            required: U256,
            /// Gas provided.
            got: U256,
        },
        /// Transaction is not valid anymore (state already has higher nonce)
        Stale,
        /// Transaction has too low fee
        /// (there is already a transaction with the same sender-nonce but higher
        /// gas price)
        TooCheapToReplace,
        /// Transaction was not imported to the queue because limit has been
        /// reached.
        LimitReached,
        /// Transaction's gas price is below threshold.
        InsufficientGasPrice {
            /// Minimal expected gas price
            minimal: U256,
            /// Transaction gas price
            got: U256,
        },
        /// Transaction's gas is below currently set minimal gas requirement.
        InsufficientGas {
            /// Minimal expected gas
            minimal: U256,
            /// Transaction gas
            got: U256,
        },
        /// Sender doesn't have enough funds to pay for this transaction
        InsufficientBalance {
            /// Senders balance
            balance: U256,
            /// Transaction cost
            cost: U256,
        },
        /// Transactions gas is higher then current gas limit
        GasLimitExceeded {
            /// Current gas limit
            limit: U256,
            /// Declared transaction gas
            got: U256,
        },
        /// Transaction's gas limit (aka gas) is invalid.
        InvalidGasLimit(OutOfBounds<U256>),
        /// Signature error
        InvalidSignature(String),
        /// Transaction too big
        TooBig,
        /// Invalid RLP encoding
        InvalidRlp(String),
        ZeroGasPrice,
        /// Transaction types have not been activated
        FutureTransactionType,
        /// Receiver with invalid type bit.
        InvalidReceiver,
        /// Transaction nonce exceeds local limit.
        TooLargeNonce,
    }
    
    impl From<keylib::Error> for TransactionError {
        fn from(err: keylib::Error) -> Self {
            TransactionError::InvalidSignature(format!("{}", err))
        }
    }
    
    impl From<rlp::DecoderError> for TransactionError {
        fn from(err: rlp::DecoderError) -> Self {
            TransactionError::InvalidRlp(format!("{}", err))
        }
    }
    
    impl fmt::Display for TransactionError {
        fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
            use self::TransactionError::*;
            let msg = match *self {
                AlreadyImported => "Already imported".into(),
                ChainIdMismatch { expected, got, space } => {
                    format!("Chain id mismatch, expected {}, got {}, space {:?}", expected, got, space)
                }
                EpochHeightOutOfBound {
                    block_height,
                    set,
                    transaction_epoch_bound,
                } => format!(
                    "EpochHeight out of bound:\
                     block_height {}, transaction epoch_height {}, transaction_epoch_bound {}",
                    block_height, set, transaction_epoch_bound
                ),
                NotEnoughBaseGas { got, required } => format!(
                    "Transaction gas {} less than intrinsic gas {}",
                    got, required
                ),
                Stale => "No longer valid".into(),
                TooCheapToReplace => "Gas price too low to replace".into(),
                LimitReached => "Transaction limit reached".into(),
                InsufficientGasPrice { minimal, got } => format!(
                    "Insufficient gas price. Min={}, Given={}",
                    minimal, got
                ),
                InsufficientGas { minimal, got } => {
                    format!("Insufficient gas. Min={}, Given={}", minimal, got)
                }
                InsufficientBalance { balance, cost } => format!(
                    "Insufficient balance for transaction. Balance={}, Cost={}",
                    balance, cost
                ),
                GasLimitExceeded { limit, got } => {
                    format!("Gas limit exceeded. Limit={}, Given={}", limit, got)
                }
                InvalidGasLimit(ref err) => format!("Invalid gas limit. {}", err),
                InvalidSignature(ref err) => {
                    format!("Transaction has invalid signature: {}.", err)
                }
                TooBig => "Transaction too big".into(),
                InvalidRlp(ref err) => {
                    format!("Transaction has invalid RLP structure: {}.", err)
                }
                ZeroGasPrice => "Zero gas price is not allowed".into(),
                FutureTransactionType => "Ethereum like transaction should have u64::MAX storage limit".into(),
                InvalidReceiver => "Sending transaction to invalid address. The first four bits of address must be 0x0, 0x1, or 0x8.".into(),
                TooLargeNonce => "Transaction nonce is too large.".into(),
            };
    
            f.write_fmt(format_args!("Transaction error ({})", msg))
        }
    }
    
    impl error::Error for TransactionError {
        fn description(&self) -> &str { "Transaction error" }
    }
    
    #[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
    pub enum Action {
        /// Create creates new contract.
        Create,
        /// Calls contract at given address.
        /// In the case of a transfer, this is the receiver's address.'
        Call(Address),
    }
    
    impl Default for Action {
        fn default() -> Action { Action::Create }
    }
    
    impl Decodable for Action {
        fn decode(rlp: &Rlp) -> Result<Self, DecoderError> {
            if rlp.is_empty() {
                Ok(Action::Create)
            } else {
                Ok(Action::Call(rlp.as_val()?))
            }
        }
    }
    
    impl Encodable for Action {
        fn rlp_append(&self, stream: &mut RlpStream) {
            match *self {
                Action::Create => stream.append_internal(&""),
                Action::Call(ref address) => stream.append_internal(address),
            };
        }
    }
    
    #[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
    pub struct AccessListItem {
        pub address: Address,
        pub storage_keys: Vec<H256>,
    }
    
    impl Encodable for AccessListItem {
        fn rlp_append(&self, s: &mut RlpStream) {
            s.begin_list(2);
            s.append(&self.address);
            s.append_list(&self.storage_keys);
        }
    }
    
    impl Decodable for AccessListItem {
        fn decode(rlp: &Rlp) -> Result<Self, DecoderError> {
            Ok(Self {
                address: rlp.val_at(0)?,
                storage_keys: rlp.list_at(1)?,
            })
        }
    }
    
    pub type AccessList = Vec<AccessListItem>;
    
    #[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
    pub enum Transaction {
        Native(TypedNativeTransaction),
        Ethereum(EthereumTransaction),
    }
    
    impl Default for Transaction {
        fn default() -> Self {
            Transaction::Native(TypedNativeTransaction::Cip155(Default::default()))
        }
    }
    
    impl From<NativeTransaction> for Transaction {
        fn from(tx: NativeTransaction) -> Self {
            Self::Native(TypedNativeTransaction::Cip155(tx))
        }
    }
    
    impl From<Eip155Transaction> for Transaction {
        fn from(tx: Eip155Transaction) -> Self {
            Self::Ethereum(EthereumTransaction::Eip155(tx))
        }
    }
    
    macro_rules! access_common_ref {
        ($field:ident, $ty:ty) => {
            pub fn $field(&self) -> &$ty {
                match self {
                    Transaction::Native(tx) => tx.$field(),
                    Transaction::Ethereum(tx) => tx.$field(),
                }
            }
        };
    }
    
    #[allow(unused)]
    macro_rules! access_common {
        ($field:ident, $ty:ident) => {
            pub fn $field(&self) -> $ty {
                match self {
                    Transaction::Native(tx) => tx.$field,
                    Transaction::Ethereum(tx) => tx.$field,
                }
            }
        };
    }
    impl Transaction {
        access_common_ref!(gas, U256);
    
        access_common_ref!(gas_price, U256);
    
        access_common_ref!(max_priority_gas_price, U256);
    
        access_common_ref!(data, Bytes);
    
        access_common_ref!(nonce, U256);
    
        access_common_ref!(action, Action);
    
        access_common_ref!(value, U256);
    
        pub fn chain_id(&self) -> Option<u32> {
            match self {
                Transaction::Native(tx) => Some(*tx.chain_id()),
                Transaction::Ethereum(tx) => tx.chain_id().clone(),
            }
        }
    
        pub fn storage_limit(&self) -> Option<u64> {
            match self {
                Transaction::Native(tx) => Some(*tx.storage_limit()),
                Transaction::Ethereum(_tx) => None,
            }
        }
    
        pub fn nonce_mut(&mut self) -> &mut U256 {
            match self {
                Transaction::Native(tx) => tx.nonce_mut(),
                Transaction::Ethereum(tx) => tx.nonce_mut(),
            }
        }
    
        pub fn type_id(&self) -> u8 {
            match self {
                Transaction::Native(TypedNativeTransaction::Cip155(_))
                | Transaction::Ethereum(EthereumTransaction::Eip155(_)) => 0,
    
                Transaction::Native(TypedNativeTransaction::Cip2930(_))
                | Transaction::Ethereum(EthereumTransaction::Eip2930(_)) => 1,
    
                Transaction::Native(TypedNativeTransaction::Cip1559(_))
                | Transaction::Ethereum(EthereumTransaction::Eip1559(_)) => 2,
            }
        }
    
        pub fn is_legacy(&self) -> bool {
            matches!(
                self,
                Transaction::Native(TypedNativeTransaction::Cip155(_))
                    | Transaction::Ethereum(EthereumTransaction::Eip155(_))
            )
        }
    
        pub fn is_2718(&self) -> bool { !self.is_legacy() }
    
        pub fn after_1559(&self) -> bool {
            matches!(
                self,
                Transaction::Native(TypedNativeTransaction::Cip1559(_))
                    | Transaction::Ethereum(EthereumTransaction::Eip1559(_))
            )
        }
    
        pub fn access_list(&self) -> Option<&AccessList> {
            match self {
                Transaction::Native(tx) => tx.access_list(),
                Transaction::Ethereum(tx) => tx.access_list(),
            }
        }
    }
    
    impl Transaction {
        pub fn priority_gas_price(&self, base_price: &U256) -> U256 {
            std::cmp::min(
                *self.max_priority_gas_price(),
                self.gas_price() - base_price,
            )
        }
    
        pub fn effective_gas_price(&self, base_price: &U256) -> U256 {
            base_price + self.priority_gas_price(base_price)
        }
    
        // This function returns the hash value used in transaction signature. It is
        // different from transaction hash. The transaction hash also contains
        // signatures.
        pub fn signature_hash(&self) -> H256 {
            let mut s = RlpStream::new();
            let mut type_prefix = vec![];
            match self {
                Transaction::Native(TypedNativeTransaction::Cip155(tx)) => {
                    s.append(tx);
                }
                Transaction::Native(TypedNativeTransaction::Cip1559(tx)) => {
                    s.append(tx);
                    type_prefix.extend_from_slice(TYPED_NATIVE_TX_PREFIX);
                    type_prefix.push(CIP1559_TYPE);
                }
                Transaction::Native(TypedNativeTransaction::Cip2930(tx)) => {
                    s.append(tx);
                    type_prefix.extend_from_slice(TYPED_NATIVE_TX_PREFIX);
                    type_prefix.push(CIP2930_TYPE);
                }
                Transaction::Ethereum(EthereumTransaction::Eip155(tx)) => {
                    s.append(tx);
                }
                Transaction::Ethereum(EthereumTransaction::Eip1559(tx)) => {
                    s.append(tx);
                    type_prefix.push(EIP1559_TYPE);
                }
                Transaction::Ethereum(EthereumTransaction::Eip2930(tx)) => {
                    s.append(tx);
                    type_prefix.push(EIP2930_TYPE);
                }
            };
            let encoded = s.as_raw();
            let mut out = vec![0; type_prefix.len() + encoded.len()];
            out[0..type_prefix.len()].copy_from_slice(&type_prefix);
            out[type_prefix.len()..].copy_from_slice(&encoded);
            keccak(&out)
        }
    
        pub fn space(&self) -> Space {
            match self {
                Transaction::Native(_) => Space::Native,
                Transaction::Ethereum(_) => Space::Ethereum,
            }
        }
    
        pub fn sign(self, secret: &Secret) -> SignedTransaction {
            let sig = ::keylib::sign(secret, &self.signature_hash())
                .expect("data is valid and context has signing capabilities; qed");
            let tx_with_sig = self.with_signature(sig);
            let public = tx_with_sig
                .recover_public()
                .expect("secret is valid so it's recoverable");
            SignedTransaction::new(public, tx_with_sig)
        }
    
        /// Signs the transaction with signature.
        pub fn with_signature(self, sig: Signature) -> TransactionWithSignature {
            TransactionWithSignature {
                transaction: TransactionWithSignatureSerializePart {
                    unsigned: self,
                    r: sig.r().into(),
                    s: sig.s().into(),
                    v: sig.v(),
                },
                hash: H256::zero(),
                rlp_size: None,
            }
            .compute_hash()
        }
    }
    
    impl MallocSizeOf for Transaction {
        fn size_of(&self, ops: &mut MallocSizeOfOps) -> usize {
            self.data().size_of(ops)
        }
    }
    
    /// Signed transaction information without verified signature.
    #[derive(Debug, Clone, Eq, PartialEq, Serialize, Deserialize)]
    pub struct TransactionWithSignatureSerializePart {
        /// Plain Transaction.
        pub unsigned: Transaction,
        /// The V field of the signature; helps describe which half of the curve
        /// our point falls in.
        pub v: u8,
        /// The R field of the signature; helps describe the point on the curve.
        pub r: U256,
        /// The S field of the signature; helps describe the point on the curve.
        pub s: U256,
    }
    
    impl Encodable for TransactionWithSignatureSerializePart {
        fn rlp_append(&self, s: &mut RlpStream) {
            match self.unsigned {
                Transaction::Native(TypedNativeTransaction::Cip155(ref tx)) => {
                    s.begin_list(4);
                    s.append(tx);
                    s.append(&self.v);
                    s.append(&self.r);
                    s.append(&self.s);
                }
                Transaction::Ethereum(EthereumTransaction::Eip155(ref tx)) => {
                    let Eip155Transaction {
                        nonce,
                        gas_price,
                        gas,
                        action,
                        value,
                        data,
                        chain_id,
                    } = tx;
                    let legacy_v = eip155_signature::add_chain_replay_protection(
                        self.v,
                        chain_id.map(|x| x as u64),
                    );
                    s.begin_list(9);
                    s.append(nonce);
                    s.append(gas_price);
                    s.append(gas);
                    s.append(action);
                    s.append(value);
                    s.append(data);
                    s.append(&legacy_v);
                    s.append(&self.r);
                    s.append(&self.s);
                }
                Transaction::Ethereum(EthereumTransaction::Eip2930(ref tx)) => {
                    s.append_raw(&[EIP2930_TYPE], 0);
                    s.begin_list(11);
                    s.append(&tx.chain_id);
                    s.append(&tx.nonce);
                    s.append(&tx.gas_price);
                    s.append(&tx.gas);
                    s.append(&tx.action);
                    s.append(&tx.value);
                    s.append(&tx.data);
                    s.append_list(&tx.access_list);
                    s.append(&self.v);
                    s.append(&self.r);
                    s.append(&self.s);
                }
                Transaction::Ethereum(EthereumTransaction::Eip1559(ref tx)) => {
                    s.append_raw(&[EIP1559_TYPE], 0);
                    s.begin_list(12);
                    s.append(&tx.chain_id);
                    s.append(&tx.nonce);
                    s.append(&tx.max_priority_fee_per_gas);
                    s.append(&tx.max_fee_per_gas);
                    s.append(&tx.gas);
                    s.append(&tx.action);
                    s.append(&tx.value);
                    s.append(&tx.data);
                    s.append_list(&tx.access_list);
                    s.append(&self.v);
                    s.append(&self.r);
                    s.append(&self.s);
                }
                Transaction::Native(TypedNativeTransaction::Cip2930(ref tx)) => {
                    s.append_raw(TYPED_NATIVE_TX_PREFIX, 0);
                    s.append_raw(&[CIP2930_TYPE], 0);
                    s.begin_list(4);
                    s.append(tx);
                    s.append(&self.v);
                    s.append(&self.r);
                    s.append(&self.s);
                }
                Transaction::Native(TypedNativeTransaction::Cip1559(ref tx)) => {
                    s.append_raw(TYPED_NATIVE_TX_PREFIX, 0);
                    s.append_raw(&[CIP1559_TYPE], 0);
                    s.begin_list(4);
                    s.append(tx);
                    s.append(&self.v);
                    s.append(&self.r);
                    s.append(&self.s);
                }
            }
        }
    }
    
    impl Decodable for TransactionWithSignatureSerializePart {
        fn decode(rlp: &Rlp) -> Result<Self, DecoderError> {
            if rlp.as_raw().len() == 0 {
                return Err(DecoderError::RlpInvalidLength);
            }
            if rlp.is_list() {
                match rlp.item_count()? {
                    4 => {
                        let unsigned: NativeTransaction = rlp.val_at(0)?;
                        let v: u8 = rlp.val_at(1)?;
                        let r: U256 = rlp.val_at(2)?;
                        let s: U256 = rlp.val_at(3)?;
                        Ok(TransactionWithSignatureSerializePart {
                            unsigned: Transaction::Native(
                                TypedNativeTransaction::Cip155(unsigned),
                            ),
                            v,
                            r,
                            s,
                        })
                    }
                    9 => {
                        let nonce: U256 = rlp.val_at(0)?;
                        let gas_price: U256 = rlp.val_at(1)?;
                        let gas: U256 = rlp.val_at(2)?;
                        let action: Action = rlp.val_at(3)?;
                        let value: U256 = rlp.val_at(4)?;
                        let data: Vec<u8> = rlp.val_at(5)?;
                        let legacy_v: u64 = rlp.val_at(6)?;
                        let r: U256 = rlp.val_at(7)?;
                        let s: U256 = rlp.val_at(8)?;
    
                        let v = eip155_signature::extract_standard_v(legacy_v);
                        let chain_id =
                            match eip155_signature::extract_chain_id_from_legacy_v(
                                legacy_v,
                            ) {
                                Some(chain_id) if chain_id > (u32::MAX as u64) => {
                                    return Err(DecoderError::Custom(
                                        "Does not support chain_id >= 2^32",
                                    ));
                                }
                                chain_id => chain_id.map(|x| x as u32),
                            };
    
                        Ok(TransactionWithSignatureSerializePart {
                            unsigned: Transaction::Ethereum(
                                EthereumTransaction::Eip155(Eip155Transaction {
                                    nonce,
                                    gas_price,
                                    gas,
                                    action,
                                    value,
                                    chain_id,
                                    data,
                                }),
                            ),
                            v,
                            r,
                            s,
                        })
                    }
                    _ => Err(DecoderError::RlpInvalidLength),
                }
            } else {
                match rlp.as_raw()[0] {
                    TYPED_NATIVE_TX_PREFIX_BYTE => {
                        if rlp.as_raw().len() <= 4
                            || rlp.as_raw()[0..3] != *TYPED_NATIVE_TX_PREFIX
                        {
                            return Err(DecoderError::RlpInvalidLength);
                        }
                        match rlp.as_raw()[3] {
                            CIP2930_TYPE => {
                                let rlp = Rlp::new(&rlp.as_raw()[4..]);
                                if rlp.item_count()? != 4 {
                                    return Err(DecoderError::RlpIncorrectListLen);
                                }
    
                                let tx = rlp.val_at(0)?;
                                let v = rlp.val_at(1)?;
                                let r = rlp.val_at(2)?;
                                let s = rlp.val_at(3)?;
                                Ok(TransactionWithSignatureSerializePart {
                                    unsigned: Transaction::Native(
                                        TypedNativeTransaction::Cip2930(tx),
                                    ),
                                    v,
                                    r,
                                    s,
                                })
                            }
                            CIP1559_TYPE => {
                                let rlp = Rlp::new(&rlp.as_raw()[4..]);
                                if rlp.item_count()? != 4 {
                                    return Err(DecoderError::RlpIncorrectListLen);
                                }
    
                                let tx = rlp.val_at(0)?;
                                let v = rlp.val_at(1)?;
                                let r = rlp.val_at(2)?;
                                let s = rlp.val_at(3)?;
                                Ok(TransactionWithSignatureSerializePart {
                                    unsigned: Transaction::Native(
                                        TypedNativeTransaction::Cip1559(tx),
                                    ),
                                    v,
                                    r,
                                    s,
                                })
                            }
                            _ => Err(DecoderError::RlpInvalidLength),
                        }
                    }
                    EIP2930_TYPE => {
                        let rlp = Rlp::new(&rlp.as_raw()[1..]);
                        if rlp.item_count()? != 11 {
                            return Err(DecoderError::RlpIncorrectListLen);
                        }
    
                        let tx = Eip2930Transaction {
                            chain_id: rlp.val_at(0)?,
                            nonce: rlp.val_at(1)?,
                            gas_price: rlp.val_at(2)?,
                            gas: rlp.val_at(3)?,
                            action: rlp.val_at(4)?,
                            value: rlp.val_at(5)?,
                            data: rlp.val_at(6)?,
                            access_list: rlp.list_at(7)?,
                        };
                        let v = rlp.val_at(8)?;
                        let r = rlp.val_at(9)?;
                        let s = rlp.val_at(10)?;
                        Ok(TransactionWithSignatureSerializePart {
                            unsigned: Transaction::Ethereum(
                                EthereumTransaction::Eip2930(tx),
                            ),
                            v,
                            r,
                            s,
                        })
                    }
                    EIP1559_TYPE => {
                        let rlp = Rlp::new(&rlp.as_raw()[1..]);
                        if rlp.item_count()? != 12 {
                            return Err(DecoderError::RlpIncorrectListLen);
                        }
    
                        let tx = Eip1559Transaction {
                            chain_id: rlp.val_at(0)?,
                            nonce: rlp.val_at(1)?,
                            max_priority_fee_per_gas: rlp.val_at(2)?,
                            max_fee_per_gas: rlp.val_at(3)?,
                            gas: rlp.val_at(4)?,
                            action: rlp.val_at(5)?,
                            value: rlp.val_at(6)?,
                            data: rlp.val_at(7)?,
                            access_list: rlp.list_at(8)?,
                        };
                        let v = rlp.val_at(9)?;
                        let r = rlp.val_at(10)?;
                        let s = rlp.val_at(11)?;
                        Ok(TransactionWithSignatureSerializePart {
                            unsigned: Transaction::Ethereum(
                                EthereumTransaction::Eip1559(tx),
                            ),
                            v,
                            r,
                            s,
                        })
                    }
                    _ => Err(DecoderError::RlpInvalidLength),
                }
            }
        }
    }
    
    impl Deref for TransactionWithSignatureSerializePart {
        type Target = Transaction;
    
        fn deref(&self) -> &Self::Target { &self.unsigned }
    }
    
    impl DerefMut for TransactionWithSignatureSerializePart {
        fn deref_mut(&mut self) -> &mut Self::Target { &mut self.unsigned }
    }
    
    /// Signed transaction information without verified signature.
    #[derive(Debug, Clone, Eq, PartialEq, Serialize, Deserialize)]
    pub struct TransactionWithSignature {
        /// Serialize part.
        pub transaction: TransactionWithSignatureSerializePart,
        /// Hash of the transaction
        #[serde(skip)]
        pub hash: H256,
        /// The transaction size when serialized in rlp
        #[serde(skip)]
        pub rlp_size: Option<usize>,
    }
    
    impl Deref for TransactionWithSignature {
        type Target = TransactionWithSignatureSerializePart;
    
        fn deref(&self) -> &Self::Target { &self.transaction }
    }
    
    impl DerefMut for TransactionWithSignature {
        fn deref_mut(&mut self) -> &mut Self::Target { &mut self.transaction }
    }
    
    impl Decodable for TransactionWithSignature {
        fn decode(tx_rlp: &Rlp) -> Result<Self, DecoderError> {
            let rlp_size = Some(tx_rlp.as_raw().len());
            // The item count of TransactionWithSignatureSerializePart is checked in
            // its decoding.
            let hash;
            let transaction;
            if tx_rlp.is_list() {
                hash = keccak(tx_rlp.as_raw());
                // Vanilla tx encoding.
                transaction = tx_rlp.as_val()?;
            } else {
                // Typed tx encoding is wrapped as an RLP string.
                let b: Vec<u8> = tx_rlp.as_val()?;
                hash = keccak(&b);
                transaction = rlp::decode(&b)?;
            };
            Ok(TransactionWithSignature {
                transaction,
                hash,
                rlp_size,
            })
        }
    }
    
    impl Encodable for TransactionWithSignature {
        fn rlp_append(&self, s: &mut RlpStream) {
            match &self.transaction.unsigned {
                Transaction::Native(TypedNativeTransaction::Cip155(_))
                | Transaction::Ethereum(EthereumTransaction::Eip155(_)) => {
                    s.append_internal(&self.transaction);
                }
                _ => {
                    // Typed tx encoding is wrapped as an RLP string.
                    s.append_internal(&rlp::encode(&self.transaction));
                }
            }
        }
    }
    
    impl TransactionWithSignature {
        pub fn new_unsigned(tx: Transaction) -> Self {
            TransactionWithSignature {
                transaction: TransactionWithSignatureSerializePart {
                    unsigned: tx,
                    s: 0.into(),
                    r: 0.into(),
                    v: 0,
                },
                hash: Default::default(),
                rlp_size: None,
            }
        }
    
        /// Used to compute hash of created transactions
        fn compute_hash(mut self) -> TransactionWithSignature {
            let hash = keccak(&*self.rlp_bytes());
            self.hash = hash;
            self
        }
    
        /// Checks whether signature is empty.
        pub fn is_unsigned(&self) -> bool { self.r.is_zero() && self.s.is_zero() }
    
        /// Construct a signature object from the sig.
        pub fn signature(&self) -> Signature {
            let r: H256 = BigEndianHash::from_uint(&self.r);
            let s: H256 = BigEndianHash::from_uint(&self.s);
            Signature::from_rsv(&r, &s, self.v)
        }
    
        /// Checks whether the signature has a low 's' value.
        pub fn check_low_s(&self) -> Result<(), keylib::Error> {
            if !self.signature().is_low_s() {
                Err(keylib::Error::InvalidSignature)
            } else {
                Ok(())
            }
        }
    
        pub fn check_y_parity(&self) -> Result<(), keylib::Error> {
            if self.is_2718() && self.v > 1 {
                // In Typed transactions (EIP-2718), v means y_parity, which must be
                // 0 or 1
                Err(keylib::Error::InvalidYParity)
            } else {
                Ok(())
            }
        }
    
        pub fn hash(&self) -> H256 { self.hash }
    
        /// Recovers the public key of the sender.
        pub fn recover_public(&self) -> Result<Public, keylib::Error> {
            Ok(recover(&self.signature(), &self.unsigned.signature_hash())?)
        }
    
        pub fn rlp_size(&self) -> usize {
            self.rlp_size.unwrap_or_else(|| self.rlp_bytes().len())
        }
    
        pub fn from_raw(raw: &[u8]) -> Result<Self, DecoderError> {
            Ok(TransactionWithSignature {
                transaction: Rlp::new(raw).as_val()?,
                hash: keccak(raw),
                rlp_size: Some(raw.len()),
            })
        }
    }
    
    impl MallocSizeOf for TransactionWithSignature {
        fn size_of(&self, ops: &mut MallocSizeOfOps) -> usize {
            self.unsigned.size_of(ops)
        }
    }
    
    /// A signed transaction with successfully recovered sender.
    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    pub struct SignedTransaction {
        pub transaction: TransactionWithSignature,
        pub sender: Address,
        pub public: Option<Public>,
    }
    
    // The default encoder for local storage.
    impl Encodable for SignedTransaction {
        fn rlp_append(&self, s: &mut RlpStream) {
            s.begin_list(3);
            s.append(&self.transaction);
            s.append(&self.sender);
            s.append(&self.public);
        }
    }
    
    impl Decodable for SignedTransaction {
        fn decode(rlp: &Rlp) -> Result<Self, DecoderError> {
            Ok(SignedTransaction {
                transaction: rlp.val_at(0)?,
                sender: rlp.val_at(1)?,
                public: rlp.val_at(2)?,
            })
        }
    }
    
    impl Deref for SignedTransaction {
        type Target = TransactionWithSignature;
    
        fn deref(&self) -> &Self::Target { &self.transaction }
    }
    
    impl DerefMut for SignedTransaction {
        fn deref_mut(&mut self) -> &mut Self::Target { &mut self.transaction }
    }
    
    impl From<SignedTransaction> for TransactionWithSignature {
        fn from(tx: SignedTransaction) -> Self { tx.transaction }
    }
    
    impl SignedTransaction {
        /// Try to verify transaction and recover sender.
        pub fn new(public: Public, transaction: TransactionWithSignature) -> Self {
            if transaction.is_unsigned() {
                SignedTransaction {
                    transaction,
                    sender: UNSIGNED_SENDER,
                    public: None,
                }
            } else {
                let sender = public_to_address(
                    &public,
                    transaction.space() == Space::Native,
                );
                SignedTransaction {
                    transaction,
                    sender,
                    public: Some(public),
                }
            }
        }
    
        pub fn new_unsigned(transaction: TransactionWithSignature) -> Self {
            SignedTransaction {
                transaction,
                sender: UNSIGNED_SENDER,
                public: None,
            }
        }
    
        pub fn set_public(&mut self, public: Public) {
            let type_nibble = self.unsigned.space() == Space::Native;
            self.sender = public_to_address(&public, type_nibble);
            self.public = Some(public);
        }
    
        /// Returns transaction sender.
        pub fn sender(&self) -> AddressWithSpace {
            self.sender.with_space(self.space())
        }
    
        pub fn nonce(&self) -> &U256 { self.transaction.nonce() }
    
        /// Checks if signature is empty.
        pub fn is_unsigned(&self) -> bool { self.transaction.is_unsigned() }
    
        pub fn hash(&self) -> H256 { self.transaction.hash() }
    
        pub fn gas(&self) -> &U256 { &self.transaction.gas() }
    
        pub fn gas_price(&self) -> &U256 { &self.transaction.gas_price() }
    
        pub fn gas_limit(&self) -> &U256 { &self.transaction.gas() }
    
        pub fn storage_limit(&self) -> Option<u64> {
            self.transaction.storage_limit()
        }
    
        pub fn rlp_size(&self) -> usize { self.transaction.rlp_size() }
    
        pub fn public(&self) -> &Option<Public> { &self.public }
    
        pub fn verify_public(&self, skip: bool) -> Result<bool, keylib::Error> {
            if self.public.is_none() {
                return Ok(false);
            }
    
            if !skip {
                let public = self.public.unwrap();
                Ok(verify_public(
                    &public,
                    &self.signature(),
                    &self.unsigned.signature_hash(),
                )?)
            } else {
                Ok(true)
            }
        }
    }
    
    impl MallocSizeOf for SignedTransaction {
        fn size_of(&self, ops: &mut MallocSizeOfOps) -> usize {
            self.transaction.size_of(ops)
        }
    }
    `
	sc := NewSouceCode(raw)
	dts, us := sc.GetDefineTypes()
	fmt.Printf("dts %+v\nus %+v\n", dts, us)
}

func TestGetDefineType(t *testing.T) {
	content := "pub type AccessList = Vec<AccessListItem>;"
	var re = regexp2.MustCompile(`^\/\/\/[^;]*?pub type (.*?) = (.*?);|pub type (.*?) = (.*?);`, regexp2.Multiline|regexp2.Singleline)
	matchs, err := re.FindStringMatch(content)
	assert.NoError(t, err)

	fmt.Printf("matchs len %d\n", len(matchs.Groups()))

	for i, item := range matchs.Groups() {
		fmt.Printf("group catpture %d: %s\n", i, item.Capture.String())
	}
}

func TestGetStructName(t *testing.T) {
	content := `
    pub struct AccessListItem {
        pub address: Address,
        pub storage_keys: Vec<H256>,
    }`
	var re = regexp2.MustCompile(`\/\/\/[^{}]+pub struct ([^\{]*) \{.*?}|pub struct ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	matchs, err := re.FindStringMatch(content)
	assert.NoError(t, err)

	fmt.Printf("matchs len %d\n", len(matchs.Groups()))
	for i, item := range matchs.Groups() {
		fmt.Printf("group catpture %d: %s\n", i, item.Capture.String())
	}
}

func TestGetSurcts(t *testing.T) {
	content := `
    pub struct AccessListItem {
        pub address: Address,
        pub storage_keys: Vec<H256>,
    }`
	var re = regexp2.MustCompile(`\/\/\/[^{}]+pub struct ([^\{]*) \{.*?}|pub struct ([^\{]*) \{.*?}`, regexp2.Multiline|regexp2.Singleline)
	sc := NewSouceCode(content)

	structs, uses := sc.getStructsOrEnums(re)
	fmt.Printf("structs %+v\nuses %+v\n", structs, uses)
}
