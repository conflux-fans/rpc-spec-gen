package rust

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/conflux-fans/rpc-spec-gen/utils"
	"github.com/dlclark/regexp2"
)

func TestParseUse(t *testing.T) {
	uses := []Use{
		// `use crate::rpc::types::eth::{BlockNumber, LocalizedTrace, TraceFilter};`,
		// `use cfx_types::H256;`,
		// `use jsonrpc_core::Result as JsonRpcResult;`,
		// `use jsonrpc_derive::rpc;`,
		// `use cfxcore::{observer::trace::Outcome, vm::CallType as CfxCallType};`,
		// `use std::{convert::{TryFrom,TryInto},convert2::{TryFrom2,TryInto2},fmt,};`,
		// `use std::{
		// 	convert::{TryFrom, TryInto},
		// 	fmt,
		// };`,
		// `use crate::{
		// 	epoch_change::Verifier,
		// 	ledger_info::{LedgerInfo, LedgerInfoWithSignatures},
		// 	on_chain_config::OnChainConfig,
		// 	validator_verifier::ValidatorVerifier,
		// };`,
		// `use crate::rpc::types::{
		// 	pos::PoSEpochReward, Account as RpcAccount, AccountPendingInfo,
		// 	AccountPendingTransactions, Block, BlockHashOrEpochNumber, Bytes,
		// 	CallRequest, CfxRpcLogFilter, CheckBalanceAgainstTransactionResponse,
		// 	EpochNumber, EstimateGasAndCollateralResponse, Log as RpcLog, PoSEconomics,
		// 	Receipt as RpcReceipt, RewardInfo as RpcRewardInfo, RpcAddress,
		// 	SponsorInfo, Status as RpcStatus, TokenSupplyInfo, Transaction,
		// };`,
		`use crate::rpc::types::{
			pos::PoSEpochReward, Account as RpcAccount, AccountPendingInfo,
			AccountPendingTransactions, Block, BlockHashOrEpochNumber, Bytes,
			CallRequest, CfxRpcLogFilter, CheckBalanceAgainstTransactionResponse,
			EpochNumber, EstimateGasAndCollateralResponse, Log as RpcLog, PoSEconomics,
			Receipt as RpcReceipt, RewardInfo as RpcRewardInfo, RpcAddress,
			SponsorInfo, Status as RpcStatus, TokenSupplyInfo, Transaction,
		};
		`,
	}

	// std, {convert::{TryFrom,TryInto},convert2::{TryFrom2,TryInto2},fmt,}
	// convert::{TryFrom,TryInto} |,convert2::{TryFrom2,TryInto2} |,fmt,

	for _, use := range uses {
		parsed := use.Parse()
		j, _ := json.MarshalIndent(parsed, "", "  ")
		fmt.Printf("%s", j)
	}
}

func TestMatchPair(t *testing.T) {
	re := regexp2.MustCompile(`\((?>[^()]+|\([^()\s]+\s(?<DEPTH>)|\)\s(?<-DEPTH>))*(?(DEPTH)(?!))\)`, regexp2.IgnoreCase)
	var str = `a()(sdf)(TOP (S (NPB (DT The) (NN question) ) (VP (VBZ remains) (SBAR-A (IN whether) (S-A (NPB (PRP they) ) (VP (MD will) (VP-A (VB be) (ADJP (JJ able) (SG (VP (TO to) (VP-A (VB help) (PUNC. us.) ) ) ) ) ) ) ) ) ) ) )
`
	m, e := re.FindStringMatch(str)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Printf("%s\n", m.String())

	for m != nil {
		m, e = re.FindNextMatch(m)
		if e != nil {
			t.Fatal(e)
		}
		if m == nil {
			break
		}
		fmt.Printf("%s\n", m.String())
	}
}

type node struct {
	name   string
	childs []node
}

func parse(str string) []node {
	if strings.TrimSpace(str) == "" {
		return nil
	}

	re := regexp2.MustCompile(`\{(?>[^{}]+|\{(?<DEPTH>)|\}(?<-DEPTH>))*(?(DEPTH)(?!))\}`, regexp2.Multiline)

	matches := []string{}
	m, e := re.FindStringMatch(str)
	if e != nil {
		panic(e)
	}
	if m == nil {
		items := strings.Split(str, ",")
		nodes := []node{}
		for _, item := range items {
			nodes = append(nodes, node{name: strings.TrimSpace(item)})
		}
		return nodes
	}

	matches = append(matches, m.String())
	for m != nil {
		m, e = re.FindNextMatch(m)
		if e != nil {
			panic(e)
		}
		if m == nil {
			break
		}
		matches = append(matches, m.String())
	}
	fmt.Printf("matches %#v\n", matches)

	replaced, e := re.Replace(str, "###", -1, -1)
	if e != nil {
		panic(e)
	}
	fmt.Printf("replaced: %s\n", replaced)

	prefixs := strings.Split(replaced, "###")
	fmt.Printf("prefixs: %s\n", prefixs)

	nodes := []node{}
	max := utils.MaxInt(len(prefixs), len(matches))
	for i := 0; i < max; i++ {

		ps := strings.Split(prefixs[i], ",")
		// fmt.Printf("ps %v: %s\n", len(ps), ps)
		for i := 0; i < len(ps); i++ {
			if strings.TrimSpace(ps[i]) == "" {
				continue
			}
			nodes = append(nodes, node{name: ps[i]})
		}

		if i < len(matches) {
			m := matches[i]
			nodes[len(nodes)-1].childs = parse(m[1 : len(m)-1])
		}

	}

	// for i, m := range matches {
	// 	ps := strings.Split(prefixs[i], ",")
	// 	fmt.Printf("ps %v: %s\n", len(ps), ps)
	// 	for i := 0; i < len(ps); i++ {
	// 		if strings.TrimSpace(ps[i]) == "" {
	// 			continue
	// 		}
	// 		nodes = append(nodes, node{name: ps[i]})
	// 	}
	// 	nodes[len(nodes)-1].childs = parse(m[1 : len(m)-1])
	// }

	// for i := len(matches) - 1; i < len(prefixs); i++ {
	// 	ps := strings.Split(prefixs[i], ",")
	// 	for i := 0; i < len(ps); i++ {
	// 		if strings.TrimSpace(ps[i]) == "" {
	// 			continue
	// 		}
	// 		nodes = append(nodes, node{name: ps[i]})
	// 	}
	// }

	return nodes
}

// func maxInt(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

func TestParseUse2(t *testing.T) {
	str := `x::{g1::{t1,t2},g2::{t3,t4},t5,g3::{sg1::{ssg1::{t6,t7}}}},t8`
	nodes := parse(str)
	fmt.Printf("%#v\n", nodes)
}
