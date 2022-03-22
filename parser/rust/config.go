package rust

import (
	"path"
	"strings"

	"github.com/Conflux-Chain/rpc-gen/config"
)

// type Config struct {
// 	RustRootPath string
// }
type RustUseTypeMeta struct {
	isBaseType bool
	isIgnore   bool
	file       string
}

func (r *RustUseTypeMeta) IsBaseType() bool {
	return r.isBaseType
}

func (r *RustUseTypeMeta) IsIgnore() bool {
	return r.isIgnore
}

func (r *RustUseTypeMeta) InFilePath() string {
	return path.Join(config.GetConfig().RustRootPath, r.file)
}

var FolderPathOfMod = map[string]string{
	"crate::rpc::types::pos": "client/src/rpc/types/pos/",
	"diem_types":             "core/src/pos/types/src/",
}

// TODO: 加一个 mod path 到 file path 的映射, 如果 RustUseTypeMetas 中找不到，从 mode path 映射中遍历解析文件查找（只需要在init时解析一遍）
// Struct exist in rust file path
var RustUseTypeMetas map[string]RustUseTypeMeta = map[string]RustUseTypeMeta{
	"cfx_types::H160":          {isBaseType: true},
	"cfx_types::H256":          {isBaseType: true},
	"cfx_types::U256":          {isBaseType: true},
	"cfx_types::U64":           {isBaseType: true},
	"H160":                     {isBaseType: true},
	"H256":                     {isBaseType: true},
	"U256":                     {isBaseType: true},
	"U64":                      {isBaseType: true},
	"address":                  {isBaseType: true},
	"String":                   {isBaseType: true},
	"bool":                     {isBaseType: true},
	"super::super::RpcAddress": {isBaseType: true},
	"u64":                      {isBaseType: true},
	"u8":                       {isBaseType: true},
	"RpcAddress":               {isBaseType: true},

	"jsonrpc_core::Result":                {isIgnore: true},
	"jsonrpc_core::ResultasJsonRpcResult": {isIgnore: true},
	"jsonrpc_derive::rpc":                 {isIgnore: true},
	"serde_derive::Serialize":             {isIgnore: true},
	"serde::de::Error":                    {isIgnore: true},
	"serde::Serialize":                    {isIgnore: true},

	"crate::rpc::types::pos::RpcTransactionStatus":     {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::RpcTransactionType":       {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::RpcTransactionPayload":    {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::RegisterPayload":          {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::ElectionPayload":          {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::UpdateVotingPowerPayload": {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::PivotDecisionPayload":     {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::RetirePayload":            {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::DisputePayload":           {file: "client/src/rpc/types/pos/transaction.rs"},
	"crate::rpc::types::pos::ConflictingVotes":         {file: "client/src/rpc/types/pos/transaction.rs"},

	"crate::rpc::types::pos::Account":        {file: "client/src/rpc/types/pos/account.rs"},
	"crate::rpc::types::pos::Block":          {file: "client/src/rpc/types/pos/block.rs"},
	"crate::rpc::types::pos::BlockNumber":    {file: "client/src/rpc/types/pos/block_number.rs"},
	"crate::rpc::types::pos::CommitteeState": {file: "client/src/rpc/types/pos/committee.rs"},
	"crate::rpc::types::pos::PoSEpochReward": {file: "client/src/rpc/types/pos/reward.rs"},
	"crate::rpc::types::pos::Status":         {file: "client/src/rpc/types/pos/status.rs"},
	"crate::rpc::types::pos::Transaction":    {file: "client/src/rpc/types/pos/transaction.rs"},

	"diem_types::epoch_state::EpochState":               {file: "core/src/pos/types/src/epoch_state.rs"},
	"diem_types::ledger_info::LedgerInfoWithSignatures": {file: "core/src/pos/types/src/ledger_info.rs"},
	"diem_types::ledger_info::LedgerInfoWithV0":         {file: "core/src/pos/types/src/ledger_info.rs"},

	"crate::rpc::types::pos::NodeLockStatus": {file: "client/src/rpc/types/pos/node_lock_status.rs"},
	"crate::rpc::types::pos::VotePowerState": {file: "client/src/rpc/types/pos/node_lock_status.rs"},

	"crate::rpc::types::pos::Signature": {file: "client/src/rpc/types/pos/block.rs"},

	"crate::rpc::types::pos::RpcCommittee":    {file: "client/src/rpc/types/pos/committee.rs"},
	"crate::rpc::types::pos::RpcTermData":     {file: "client/src/rpc/types/pos/committee.rs"},
	"crate::rpc::types::pos::NodeVotingPower": {file: "client/src/rpc/types/pos/committee.rs"},

	"crate::rpc::types::pos::Decision": {file: "client/src/rpc/types/pos/decision.rs"},

	"crate::rpc::types::pos::Reward": {file: "client/src/rpc/types/pos/reward.rs"},

	"super::Decision": {file: "client/src/rpc/types/pos/decision.rs"},
	"diem_types::block_info::PivotBlockDecision": {file: "core/src/pos/types/src/block_info.rs"},

	"crate::validator_verifier::ValidatorVerifier": {file: "core/src/pos/types/src/validator_verifier.rs"},
}

func GetUseTypeMeta(useType UseType) (*RustUseTypeMeta, bool) {
	v, ok := RustUseTypeMetas[useType.String()]
	return &v, ok
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
