package specconfig

import "github.com/conflux-fans/rpc-spec-gen/code_gen/openrpc/types"

var cfxServerConfig = []*types.Server{
	{
		Name:    "conflux core space mainnet RPC",
		URL:     "https://main.confluxrpc.com",
		Summary: "The mainnet RPC server for Conflux core space. chainId: 1029",
	},
	{
		Name:    "conflux core space mainnet websocket RPC",
		URL:     "wss://main.confluxrpc.com/ws",
		Summary: "The mainnet RPC server with websocket protocal for Conflux core space. chainId: 1029",
	},
	{
		Name:    "conflux core space testnet RPC",
		URL:     "https://test.confluxrpc.com",
		Summary: "The testnet RPC server for Conflux core space. chainId: 1",
	},
	{
		Name:    "conflux core space testnet websocket RPC",
		URL:     "wss://test.confluxrpc.com/ws",
		Summary: "The testnet RPC server with websocket protocal for Conflux core space. chainId: 1",
	},
}

// Mainnet	1029	https://confluxscan.io	https://main.confluxrpc.com
// wss://main.confluxrpc.com/ws
// Testnet	1	https://testnet.confluxscan.io	https://test.confluxrpc.com
// wss://test.confluxrpc.com/ws

var cfxInfo = types.Info{
	Version:     "0.1.0",
	Description: "A specification of the standard interface of Conflux clients.",
	Title:       "Conflux JSON-RPC Specification",
	License: &types.License{
		Name: "CC0-1.0",
		URL:  "https://creativecommons.org/publicdomain/zero/1.0/legalcode",
	},
}

func getCfxSpaceSpecConfig() SpecConfig {
	s := SpecConfig{}
	s.Space = "cfx_space"
	// s.TraitName = "Eth"
	s.Info = &cfxInfo
	// s.Methods = ethMethodConfigs
	s.Servers = cfxServerConfig
	return s
}
