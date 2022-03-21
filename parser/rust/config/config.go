package config

import (
	"path"
	"strings"

	mainconfig "github.com/Conflux-Chain/rpc-gen/config"
	"github.com/Conflux-Chain/rpc-gen/parser/rust"
)

type Config struct {
	RustRootPath string
}
type RustUseTypeMeta struct {
	isBaseType bool
	file       string
}

func (r *RustUseTypeMeta) IsBaseType() bool {
	return r.isBaseType
}

func (r *RustUseTypeMeta) InFilePath() string {
	return path.Join(mainconfig.GetConfig().RustRootPath, r.file)
}

// Struct exist in rust file path
var RustUseTypeMetas map[string]RustUseTypeMeta = map[string]RustUseTypeMeta{
	"cfx_types::H160":          {true, ""},
	"cfx_types::H256":          {true, ""},
	"cfx_types::U256":          {true, ""},
	"cfx_types::U64":           {true, ""},
	"H160":                     {true, ""},
	"H256":                     {true, ""},
	"U256":                     {true, ""},
	"U64":                      {true, ""},
	"address":                  {true, ""},
	"String":                   {true, ""},
	"bool":                     {true, ""},
	"super::super::RpcAddress": {true, ""},
	"u64":                      {true, ""},
	"u8":                       {true, ""},
	"RpcAddress":               {true, ""},

	// FIXME: remove this
	// "crate::rpc::types::pos::RpcTransactionPayload": {true, ""},
	// "crate::rpc::types::pos::RpcTransactionStatus":  {true, ""},
	// "crate::rpc::types::pos::RpcTransactionType":    {true, ""},

	"crate::rpc::types::pos::RpcTransactionStatus":     {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::RpcTransactionType":       {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::RpcTransactionPayload":    {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::RegisterPayload":          {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::ElectionPayload":          {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::UpdateVotingPowerPayload": {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::PivotDecisionPayload":     {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::RetirePayload":            {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::DisputePayload":           {false, "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::ConflictingVotes":         {false, "client/src/rpc/types/pos/transaction.rs"},

	"crate::rpc::types::pos::Account":        {false, "client/src/rpc/types/pos/account.rs"},
	"crate::rpc::types::pos::Block":          {false, "client/src/rpc/types/pos/block.rs"},
	"crate::rpc::types::pos::BlockNumber":    {false, "client/src/rpc/types/pos/block_number.rs"},
	"crate::rpc::types::pos::CommitteeState": {false, "client/src/rpc/types/pos/committee.rs"},
	"crate::rpc::types::pos::PoSEpochReward": {false, "client/src/rpc/types/pos/reward.rs"},
	"crate::rpc::types::pos::Status":         {false, "client/src/rpc/types/pos/status.rs"},
	"crate::rpc::types::pos::Transaction":    {false, "client/src/rpc/types/pos/transaction.rs"},

	"diem_types::epoch_state::EpochState":               {false, "core/src/pos/types/src/epoch_state.rs"},
	"diem_types::ledger_info::LedgerInfoWithSignatures": {false, "core/src/pos/types/src/ledger_info.rs"},
	"diem_types::ledger_info::LedgerInfoWithV0":         {false, "core/src/pos/types/src/ledger_info.rs"},

	"crate::rpc::types::pos::NodeLockStatus": {false, "client/src/rpc/types/pos/node_lock_status.rs"},
	"crate::rpc::types::pos::VotePowerState": {false, "client/src/rpc/types/pos/node_lock_status.rs"},

	"crate::rpc::types::pos::Signature": {false, "client/src/rpc/types/pos/block.rs"},

	"crate::rpc::types::pos::RpcCommittee":    {false, "client/src/rpc/types/pos/committee.rs"},
	"crate::rpc::types::pos::RpcTermData":     {false, "client/src/rpc/types/pos/committee.rs"},
	"crate::rpc::types::pos::NodeVotingPower": {false, "client/src/rpc/types/pos/committee.rs"},

	"crate::rpc::types::pos::Decision": {false, "client/src/rpc/types/pos/decision.rs"},

	"crate::rpc::types::pos::Reward": {false, "client/src/rpc/types/pos/reward.rs"},

	"super::Decision": {false, "client/src/rpc/types/pos/decision.rs"},
	"diem_types::block_info::PivotBlockDecision": {false, "core/src/pos/types/src/block_info.rs"},

	"crate::validator_verifier::ValidatorVerifier": {false, "core/src/pos/types/src/validator_verifier.rs"},
}

var IgnoredUseTypes map[string]bool = map[string]bool{
	"jsonrpc_core::Result":                true,
	"jsonrpc_core::ResultasJsonRpcResult": true,
	"jsonrpc_derive::rpc":                 true,
	"serde_derive::Serialize":             true,
	"serde::de::Error":                    true,
	"serde::Serialize":                    true,
}

func GetUseTypeMeta(useType rust.UseType) *RustUseTypeMeta {
	v, ok := RustUseTypeMetas[useType.String()]
	if !ok {
		return nil
	}
	return &v
}

// check if type is base type by fullname or name
// for example: both `cfx_types::H160` and `H160` are base type
func IsBaseType(typeName string) bool {
	if m, ok := RustUseTypeMetas[typeName]; ok {
		return m.isBaseType
	}

	for k, v := range RustUseTypeMetas {
		if strings.HasSuffix(k, typeName) {
			return v.isBaseType
		}
	}
	return false
}
