package rust

import (
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/conflux-fans/rpc-spec-gen/config"
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

// mod path 到 file path 的映射, 如果 RustUseTypeMetas 中找不到，从 mode path 映射中遍历解析文件查找（只需要在init时解析一遍）
var FolderPathOfMod = map[string]string{
	"crate::rpc::types":         "client/src/rpc/types/",
	"crate::rpc::types::pos":    "client/src/rpc/types/pos/",
	"diem_types":                "core/src/pos/types/src/",
	"cfxcore::transaction_pool": "core/src/transaction_pool/",
	"primitives":                "primitives/src/",

	"crate::rpc::types::eth": "client/src/rpc/types/eth/",
}

// Struct exist in rust file path
var RustUseTypeMetas map[string]RustUseTypeMeta = map[string]RustUseTypeMeta{
	// "cfx_types::H160":          {isBaseType: true},
	// "cfx_types::H256":          {isBaseType: true},
	// "cfx_types::U256":          {isBaseType: true},
	// "cfx_types::U64":           {isBaseType: true},
	"H64":         {isBaseType: true},
	"H160":        {isBaseType: true},
	"H256":        {isBaseType: true},
	"H2048":       {isBaseType: true},
	"U256":        {isBaseType: true},
	"U64":         {isBaseType: true},
	"address":     {isBaseType: true},
	"String":      {isBaseType: true},
	"bool":        {isBaseType: true},
	"u64":         {isBaseType: true},
	"u8":          {isBaseType: true},
	"usize":       {isBaseType: true},
	"RpcAddress":  {isBaseType: true},
	"Address":     {isBaseType: true},
	"Network":     {isBaseType: true},
	"Bytes":       {isBaseType: true},
	"PosBlockId":  {isBaseType: true},
	"Bloom":       {isBaseType: true},
	"StorageRoot": {isBaseType: true},

	"jsonrpc_core::Result":                {isIgnore: true},
	"jsonrpc_core::ResultasJsonRpcResult": {isIgnore: true},
	"jsonrpc_derive::rpc":                 {isIgnore: true},
	"serde_derive::Serialize":             {isIgnore: true},
	"serde::de::Error":                    {isIgnore: true},
	"serde::Serialize":                    {isIgnore: true},
	"jsonrpc_core::BoxFuture":             {isIgnore: true},

	// "crate::rpc::types::pos::RpcTransactionStatus":     {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::RpcTransactionType":       {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::RpcTransactionPayload":    {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::RegisterPayload":          {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::ElectionPayload":          {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::UpdateVotingPowerPayload": {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::PivotDecisionPayload":     {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::RetirePayload":            {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::DisputePayload":           {file: "client/src/rpc/types/pos/transaction.rs"},
	// "crate::rpc::types::pos::ConflictingVotes":         {file: "client/src/rpc/types/pos/transaction.rs"},

	// "crate::rpc::types::pos::Account":        {file: "client/src/rpc/types/pos/account.rs"},
	// "crate::rpc::types::pos::Block":          {file: "client/src/rpc/types/pos/block.rs"},
	// "crate::rpc::types::pos::BlockNumber":    {file: "client/src/rpc/types/pos/block_number.rs"},
	// "crate::rpc::types::pos::CommitteeState": {file: "client/src/rpc/types/pos/committee.rs"},
	// "crate::rpc::types::pos::PoSEpochReward": {file: "client/src/rpc/types/pos/reward.rs"},
	// "crate::rpc::types::pos::Status":         {file: "client/src/rpc/types/pos/status.rs"},
	// "crate::rpc::types::pos::Transaction":    {file: "client/src/rpc/types/pos/transaction.rs"},

	// "crate::rpc::types::pos::NodeLockStatus": {file: "client/src/rpc/types/pos/node_lock_status.rs"},
	// "crate::rpc::types::pos::VotePowerState": {file: "client/src/rpc/types/pos/node_lock_status.rs"},

	// "crate::rpc::types::pos::Signature": {file: "client/src/rpc/types/pos/block.rs"},

	// "crate::rpc::types::pos::RpcCommittee":    {file: "client/src/rpc/types/pos/committee.rs"},
	// "crate::rpc::types::pos::RpcTermData":     {file: "client/src/rpc/types/pos/committee.rs"},
	// "crate::rpc::types::pos::NodeVotingPower": {file: "client/src/rpc/types/pos/committee.rs"},

	// "crate::rpc::types::pos::Decision": {file: "client/src/rpc/types/pos/decision.rs"},

	// "crate::rpc::types::pos::Reward": {file: "client/src/rpc/types/pos/reward.rs"},

	// "diem_types::epoch_state::EpochState":               {file: "core/src/pos/types/src/epoch_state.rs"},
	// "diem_types::ledger_info::LedgerInfoWithSignatures": {file: "core/src/pos/types/src/ledger_info.rs"},
	// "diem_types::ledger_info::LedgerInfoWithV0":         {file: "core/src/pos/types/src/ledger_info.rs"},
	// "diem_types::block_info::PivotBlockDecision":   {file: "core/src/pos/types/src/block_info.rs"},

	"super::Decision": {file: "client/src/rpc/types/pos/decision.rs"},
	"crate::validator_verifier::ValidatorVerifier": {file: "core/src/pos/types/src/validator_verifier.rs"},

	"super::Transaction": {file: "client/src/rpc/types/transaction.rs"},
	"cfx_types::Space":   {file: "cfx_types/src/lib.rs"},
	// "cfxcore::transaction_pool::TransactionStatus": {file: "core/src/transaction_pool/transaction_pool_inner.rs"},
	"super::BlockTransactions": {file: "client/src/rpc/types/block.rs"},
	"super::EpochNumber":       {file: "client/src/rpc/types/epoch_number.rs"},
}

func GetUseTypeMeta(useType UseType) (*RustUseTypeMeta, bool) {

	if useType.Name == "RpcAddress" {
		time.Sleep(0)
	}

	if v, ok := RustUseTypeMetas[useType.String()]; ok {
		return &v, true
	}

	if IsBaseType(useType.Name) {
		return &RustUseTypeMeta{isBaseType: true}, true
	}

	// 根据modpath解析folder并填充到RustUseTypeMetas中
	if f, ok := getFolderPath(useType); ok {

		folderPath := path.Join(config.GetConfig().RustRootPath, f)
		// 遍历解析文件夹下的文件
		files, err := ioutil.ReadDir(folderPath)
		if err != nil {
			logger.WithField("folderPath", folderPath).WithError(err).Panic("read dir error")
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasSuffix(file.Name(), ".rs") {
				// 填充RustUseTypeMetas
				filePath := path.Join(folderPath, file.Name())
				_content, err := ioutil.ReadFile(filePath)
				content := SourceCode(_content)
				if err != nil {
					logger.WithField("filePath", filePath).WithError(err).Panic("read file error")
				}
				structs, _ := content.GetStructs()
				enums, _ := content.GetEnums()
				defineTypes, _ := content.GetDefineTypes()
				for s := range structs {
					sUsetype := UseType{useType.ModPath, s, s}
					meta := RustUseTypeMeta{
						file: strings.TrimPrefix(filePath, config.GetConfig().RustRootPath),
					}
					addTypeMetaIfNotExist(sUsetype, meta)
				}
				for e := range enums {
					eUsetype := UseType{useType.ModPath, e, e}
					meta := RustUseTypeMeta{
						file: strings.TrimPrefix(filePath, config.GetConfig().RustRootPath),
					}
					addTypeMetaIfNotExist(eUsetype, meta)
				}
				for d := range defineTypes {
					dUsetype := UseType{useType.ModPath, d, d}
					meta := RustUseTypeMeta{
						file: strings.TrimPrefix(filePath, config.GetConfig().RustRootPath),
					}
					addTypeMetaIfNotExist(dUsetype, meta)
				}
			}
		}
		if v, ok := RustUseTypeMetas[useType.String()]; ok {
			return &v, true
		}
	}
	return nil, false
}

func getFolderPath(useType UseType) (string, bool) {
	if f, ok := FolderPathOfMod[useType.ModPathString()]; ok {
		return f, ok
	}

	for k, v := range FolderPathOfMod {
		if strings.HasPrefix(useType.ModPathString(), k) {
			return v, true
		}
	}

	return "", false
}

func addTypeMetaIfNotExist(useType UseType, meta RustUseTypeMeta) {
	if useType.IsBaseType() {
		return
	}

	if _, ok := RustUseTypeMetas[useType.String()]; !ok {
		RustUseTypeMetas[useType.String()] = meta
	}
}

// check if type is base type by fullname or name
// for example: both `cfx_types::H160` and `H160` are base type
func IsBaseType(typeName string) bool {
	if m, ok := RustUseTypeMetas[typeName]; ok {
		return m.isBaseType
	}

	for k, v := range RustUseTypeMetas {
		if strings.HasSuffix(typeName, k) {
			return v.isBaseType
		}
	}
	return false
}
