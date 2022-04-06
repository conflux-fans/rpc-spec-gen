package specconfig

import "github.com/conflux-fans/rpc-spec-gen/code_gen/openrpc/types"

var ethMethodConfigs map[string]*MethodConfig = map[string]*MethodConfig{
	"web3_clientVersion": {
		Summary:     "current client version",
		ResultName:  "clientVersion",
		Description: "client version",
	},
	"web3_sha3": {
		Summary:     "Hashes data",
		ParamNames:  []string{"data"},
		ResultName:  "hashedData",
		Description: "Keccak-256 hash of the given data",
	},
	"net_listening": {
		Summary:     "returns listening status",
		ResultName:  "netListeningResult",
		Description: "`true` if listening is active or `false` if listening is not active",
	},
	"net_peerCount": {
		Summary:     "number of peers",
		ResultName:  "quantity",
		Description: "number of connected peers.",
	},
	"net_version": {
		Summary:     "Network identifier associated with network",
		ResultName:  "networkId",
		Description: "Network ID associated with the current network",
	},
	"eth_blockNumber": {
		Summary: "Returns the number of most recent block.",
	},
	"eth_call": {
		Summary:     "Executes a new message call (locally) immediately without creating a transaction on the block chain.",
		ParamNames:  []string{"transaction", "blocknumber"},
		ResultName:  "returnValue",
		Description: "The return value of the executed contract",
	},
	"eth_chainId": {
		Summary:     "Returns the currently configured chain id",
		ResultName:  "chainId",
		Description: "hex format integer of the current chain id. Defaults are ETC=61, ETH=1, Morden=62.",
	},
	"eth_coinbase": {
		Summary:     "Returns the client coinbase address.",
		ResultName:  "address",
		Description: "The address owned by the client that is used as default for things like the mining reward",
	},
	"eth_estimateGas": {
		Summary:     "Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain. Note that the estimate may be significantly more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.",
		ResultName:  "gasUsed",
		Description: "The amount of gas used",
	},
	"eth_gasPrice": {
		Summary: "Returns the current price per gas in wei",
	},
	"eth_getBalance": {
		Summary:    "Returns Ether balance of a given or account or contract",
		ParamNames: []string{"address", "blockNumber"},
		ResultName: "getBalanceResult",
	},
	"eth_getBlockByHash": {
		Summary:    "Gets a block for a given hash",
		ParamNames: []string{"blockHash", "includeTransactions"},
		ResultName: "getBlockByHashResult",
	},
	"eth_getBlockByNumber": {
		Summary:    "Gets a block for a given number",
		ParamNames: []string{"blocknumber", "includeTransactions"},
		ResultName: "getBlockByNumberResult",
	},
	"eth_getBlockTransactionCountByHash": {
		Summary:     "Returns the number of transactions in a block from a block matching the given block hash.",
		ParamNames:  []string{"blockhash"},
		ResultName:  "blockTransactionCountByHash",
		Description: "The Number of total transactions in the given block",
	},
	"eth_getBlockTransactionCountByNumber": {
		Summary:     "Returns the number of transactions in a block from a block matching the given block number.",
		ParamNames:  []string{"blocknumber"},
		ResultName:  "blockTransactionCountByHash",
		Description: "The Number of total transactions in the given block",
	},
	"eth_getCode": {
		Summary:    "Returns code at a given contract address",
		ParamNames: []string{"address", "blockNumber"},
		ResultName: "bytes",
	},
	"eth_getFilterChanges": {
		Summary:    "Polling method for a filter, which returns an array of logs which occurred since last poll.",
		ParamNames: []string{"filterId"},
		ResultName: "logResult",
	},
	"eth_getFilterLogs": {
		Summary:    "Returns an array of all logs matching filter with given id.",
		ParamNames: []string{"filterId"},
	},
	"eth_getRawTransactionByHash": {
		Summary:     "Returns raw transaction data of a transaction with the given hash.",
		ParamNames:  []string{"transactionhash"},
		ResultName:  "rawTransactionByHash",
		Description: "The raw transaction data",
	},
	"eth_getRawTransactionByBlockHashAndIndex": {
		Summary:     "Returns raw transaction data of a transaction with the block hash and index of which it was mined.",
		ParamNames:  []string{"blockhash", "index"},
		ResultName:  "rawTransaction",
		Description: "The raw transaction data",
	},
	"eth_getRawTransactionByBlockNumberAndIndex": {
		Summary:     "Returns raw transaction data of a transaction with the block number and index of which it was mined.",
		ParamNames:  []string{"blocknumber", "index"},
		ResultName:  "rawTransaction",
		Description: "The raw transaction data",
	},
	"eth_getLogs": {
		Summary:    "Returns an array of all logs matching a given filter object.",
		ParamNames: []string{"filter"},
	},
	"eth_getStorageAt": {
		Summary:    "Gets a storage value from a contract address, a position, and an optional blockNumber",
		ParamNames: []string{"address", "position", "blocknumber"},
		ResultName: "dataWord",
	},
	"eth_getTransactionByBlockHashAndIndex": {
		Summary:    "Returns the information about a transaction requested by the block hash and index of which it was mined.",
		ParamNames: []string{"blockhash", "index"},
	},
	"eth_getTransactionByBlockNumberAndIndex": {
		Summary:    "Returns the information about a transaction requested by the block number and index of which it was mined.",
		ParamNames: []string{"blocknumber", "index"},
	},
	"eth_getTransactionByHash": {
		Summary:    "Returns the information about a transaction requested by transaction hash.",
		ParamNames: []string{"transactionhash"},
	},
	"eth_getTransactionCount": {
		Summary:    "Returns the number of transactions sent from an address",
		ParamNames: []string{"address", "blocknumber"},
		ResultName: "transactionCount",
	},
	"eth_getTransactionReceipt": {
		Summary:     "Returns the receipt information of a transaction by its hash.",
		ParamNames:  []string{"transactionhash"},
		ResultName:  "transactionReceiptResult",
		Description: "returns either a receipt or null",
	},
	"eth_getUncleByBlockHashAndIndex": {
		Summary:    "Returns information about a uncle of a block by hash and uncle index position.",
		ParamNames: []string{"blockhash", "index"},
		ResultName: "uncle",
	},
	"eth_getUncleByBlockNumberAndIndex": {
		Summary:     "Returns information about a uncle of a block by hash and uncle index position.",
		ParamNames:  []string{"uncleBlockNumber", "index"},
		ResultName:  "uncleResult",
		Description: "returns an uncle block or null",
	},
	"eth_getUncleCountByBlockHash": {
		Summary:     "Returns the number of uncles in a block from a block matching the given block hash.",
		ParamNames:  []string{"blockhash"},
		ResultName:  "uncleCountResult",
		Description: "The Number of total uncles in the given block",
	},
	"eth_getUncleCountByBlockNumber": {
		Summary:    "Returns the number of uncles in a block from a block matching the given block number.",
		ParamNames: []string{"blocknumber"},
	},
	"eth_getProof": {
		Summary:    "Returns the account- and storage-values of the specified account including the Merkle-proof.",
		ParamNames: []string{"address", "storageKeys", "blocknumber"},
		ResultName: "account",
	},
	"eth_getWork": {
		Summary:    "Returns the hash of the current block, the seedHash, and the boundary condition to be met ('target').",
		ResultName: "work",
	},
	"eth_hashrate": {
		Summary:     "Returns the number of hashes per second that the node is mining with.",
		ResultName:  "hashesPerSecond",
		Description: "Integer of the number of hashes per second",
	},
	"eth_mining": {
		Summary:     "Returns true if client is actively mining new blocks.",
		ResultName:  "mining",
		Description: "Whether or not the client is mining",
	},
	"eth_newBlockFilter": {
		Summary: "Creates a filter in the node, to notify when a new block arrives. To check if the state has changed, call eth_getFilterChanges.",
	},
	"eth_newFilter": {
		Summary:     "Creates a filter object, based on filter options, to notify when the state changes (logs). To check if the state has changed, call eth_getFilterChanges.",
		ParamNames:  []string{"filter"},
		ResultName:  "filterId",
		Description: "The filter ID for use in `eth_getFilterChanges`",
	},
	"eth_newPendingTransactionFilter": {
		Summary: "Creates a filter in the node, to notify when new pending transactions arrive. To check if the state has changed, call eth_getFilterChanges.",
	},
	"eth_pendingTransactions": {
		Summary:    "Returns the transactions that are pending in the transaction pool and have a from address that is one of the accounts this node manages.",
		ResultName: "pendingTransactions",
	},
	"eth_protocolVersion": {
		Summary:     "Returns the current ethereum protocol version.",
		ResultName:  "protocolVersion",
		Description: "The current ethereum protocol version",
	},
	"eth_sendRawTransaction": {
		Summary:     "Creates new message call transaction or a contract creation for signed transactions.",
		ParamNames:  []string{"signedTransactionData"},
		ResultName:  "transactionHash",
		Description: "The transaction hash, or the zero hash if the transaction is not yet available.",
	},
	"eth_submitHashrate": {
		Summary:     "Used for submitting mining hashrate.",
		ParamNames:  []string{"hashRate", "id"},
		ResultName:  "submitHashRateSuccess",
		Description: "whether of not submitting went through successfully",
	},
	"eth_submitWork": {
		Summary:     "Used for submitting a proof-of-work solution.",
		ParamNames:  []string{"nonce", "powHash", "mixHash"},
		ResultName:  "solutionValid",
		Description: "returns true if the provided solution is valid, otherwise false.",
	},
	"eth_syncing": {
		Summary:    "Returns an object with data about the sync status or false.",
		ResultName: "syncing",
	},
	"eth_uninstallFilter": {
		Summary:     "Uninstalls a filter with given id. Should always be called when watch is no longer needed. Additionally Filters timeout when they aren't requested with eth_getFilterChanges for a period of time.",
		ParamNames:  []string{"filterId"},
		ResultName:  "filterUninstalledSuccess",
		Description: "returns true if the filter was successfully uninstalled, false otherwise.",
	},
}

var ethServerConfig = []*types.Server{
	{
		Name:    "eSpace mainnet RPC",
		URL:     "https://evm.confluxrpc.com",
		Summary: "The mainnet RPC server for Conflux eSpace. chainId: 1030",
	},
	{
		Name:    "eSpace testnet RPC",
		URL:     "https://evmtestnet.confluxrpc.com",
		Summary: "The testnet RPC server for Conflux eSpace. chainId: 71",
	},
}

var ethInfo = types.Info{
	Version:     "0.1.0",
	Description: "A specification of the standard interface of Ethereum clients.",
	Title:       "Ethereum JSON-RPC Specification",
	License: &types.License{
		Name: "CC0-1.0",
		URL:  "https://creativecommons.org/publicdomain/zero/1.0/legalcode",
	},
}

func getEthSpaceSpecConfig() SpecConfig {
	s := SpecConfig{}
	s.Space = "eth_space"
	// s.TraitName = "Eth"
	s.Info = &ethInfo
	s.Methods = ethMethodConfigs
	s.Servers = ethServerConfig
	return s
}
