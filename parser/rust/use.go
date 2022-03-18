package rust

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Conflux-Chain/rpc-gen/utils"
	"github.com/dlclark/regexp2"
)

type Use string

type UseType struct {
	ModPath []string
	Alias   string
	Name    string
}

type useBody string
type useItem string
type useNode struct {
	name   string
	childs []useNode
}

var (
	re_PAIR  = regexp2.MustCompile(`\{(?>[^{}]+|\{(?<DEPTH>)|\}(?<-DEPTH>))*(?(DEPTH)(?!))\}`, regexp2.Multiline)
	re_ALIAS = regexp.MustCompile(`(.*) as (.*)`)
)

func MustNewUseType(usetype string) UseType {
	u := useItem(usetype)
	return u.Parse()
	// panic("invalid use type")
}

func (r *Use) Parse() []UseType {
	body := r.body()
	nodes := body.toNodes()
	// fmt.Printf("nodes %#v\n", nodes)

	flattens := []string{}
	for _, node := range nodes {
		flattens = append(flattens, node.flatten()...)
	}

	rts := []UseType{}
	for i := range flattens {
		rts = append(rts, useItem(flattens[i]).Parse())
	}
	return rts
	// return UseParsed{
	// 	rts,
	// }
}

func (r UseType) String() string {
	modPath := strings.Join(r.ModPath, "::")
	if modPath == "" {
		return r.Name
	}
	return fmt.Sprintf("%s::%s", modPath, r.Name)
}

func (r *Use) body() useBody {
	useReg := regexp.MustCompile(`(?mUs)use (.*);`)
	useFinded := useReg.FindStringSubmatch(string(*r))
	// fmt.Printf("useFinded %v\n", useFinded)
	if len(useFinded) == 0 {
		logger.WithField("use", *r).Panic("not found use body")
	}
	return useBody(useFinded[1])
}

func (b *useBody) toNodes() []useNode {
	// trimd := regexp.MustCompile(`\s+`).ReplaceAllString(string(ru), "")
	str := strings.TrimSpace(string(*b))
	if str == "" {
		return nil
	}

	// re := regexp2.MustCompile(`\{(?>[^{}]+|\{(?<DEPTH>)|\}(?<-DEPTH>))*(?(DEPTH)(?!))\}`, regexp2.Multiline)

	matches := []string{}
	m, e := re_PAIR.FindStringMatch(str)
	if e != nil {
		panic(e)
	}
	if m == nil {
		items := strings.Split(str, ",")
		nodes := []useNode{}
		for _, item := range items {
			if strings.TrimSpace(item) == "" {
				continue
			}
			nodes = append(nodes, useNode{name: strings.TrimSpace(item)})
		}
		return nodes
	}

	matches = append(matches, m.String())
	for m != nil {
		m, e = re_PAIR.FindNextMatch(m)
		if e != nil {
			panic(e)
		}
		if m == nil {
			break
		}
		matches = append(matches, m.String())
	}
	// fmt.Printf("matches %#v\n", matches)

	replaced, e := re_PAIR.Replace(str, "###", -1, -1)
	if e != nil {
		panic(e)
	}
	// fmt.Printf("replaced: %s\n", replaced)

	prefixs := strings.Split(replaced, "###")
	// fmt.Printf("prefixs: %s\n", prefixs)

	nodes := []useNode{}
	max := utils.MaxInt(len(prefixs), len(matches))
	for i := 0; i < max; i++ {

		ps := strings.Split(prefixs[i], ",")
		fmt.Printf("ps %v: %s\n", len(ps), ps)
		for i := 0; i < len(ps); i++ {
			if strings.TrimSpace(ps[i]) == "" {
				continue
			}
			nodes = append(nodes, useNode{name: ps[i]})
		}

		if i < len(matches) {
			m := matches[i]
			subBody := useBody(m[1 : len(m)-1])
			nodes[len(nodes)-1].childs = subBody.toNodes()
		}

	}
	logger.WithField("nodes", nodes).WithField("use", *b).Debug("use to nodes")
	return nodes
}

func (u *useNode) flatten() []string {

	if len(u.childs) == 0 {
		return []string{u.name}
	}

	flattened := []string{}
	for _, child := range u.childs {
		fs := child.flatten()
		for _, f := range fs {
			flattened = append(flattened, u.name+f)
		}
	}
	return flattened
}

// // use crate::rpc::types::eth::{BlockNumber, LocalizedTrace, TraceFilter};
// // use cfx_types::H256;
// // use jsonrpc_core::Result as JsonRpcResult;
// // use jsonrpc_derive::rpc;
// func (ru RustUse) Parse() RustUseParsed {
// 	result := RustUseParsed{}

// 	trimd := regexp.MustCompile(`\s+`).ReplaceAllString(string(ru), "")
// 	finds := regexp.MustCompile(`(?U)(.*)::{(.*)}`).FindStringSubmatch(trimd)
// 	if finds == nil {
// 		result.Types = append(result.Types, rustUseItem(ru).Parse())
// 		return result
// 	}

// 	modPath, itemsJoined := finds[1], finds[2]

// 	// 如果 itemsJoined 还是 {...},

// 	items := strings.Split(itemsJoined, ",")
// 	for _, item := range items {
// 		item = strings.TrimSpace(item)
// 		if item == "" {
// 			continue
// 		}
// 		result.Types = append(result.Types, RustUse(modPath+"::"+item).Parse().Types...)
// 	}
// 	return result
// }

// func (ru RustUse) Flatten() []rustUseItem {
// 	result := []rustUseItem{}
// 	trimd := regexp.MustCompile(`\s+`).ReplaceAllString(string(ru), "")
// 	finds := regexp.MustCompile(`(?U)(.*)::{(.*)}`).FindStringSubmatch(trimd)
// 	if finds == nil {
// 		result = append(result, rustUseItem(ru))
// 		return result
// 	}

// 	// finds[1]: std, finds[2]: convert::{TryFrom,TryInto},convert2::{TryFrom2,TryInto2},fmt,
// 	modPath, itemsJoined := finds[1], finds[2]
// 	subRustUses := RustUse(itemsJoined).Flatten()
// 	for _, subUse := range subRustUses {
// 		result = append(result, rustUseItem(modPath+"::"+string(subUse)))
// 	}
// 	return result
// }

func (ri useItem) Parse() UseType {
	parts := strings.Split(string(ri), "::")
	last := parts[len(parts)-1]

	//jsonrpc_core::Result as JsonRpcResult
	finds := re_ALIAS.FindStringSubmatch(last)
	if len(finds) > 0 {
		return UseType{
			ModPath: parts[0 : len(parts)-1],
			Name:    strings.TrimSpace(finds[1]),
			Alias:   strings.TrimSpace(finds[2]),
		}
	}
	// // cfx_types::H256
	// matched, _ := regexp.MatchString(`^\w+$`, last)
	// if matched {
	// }
	return UseType{
		ModPath: parts[0 : len(parts)-1],
		Name:    strings.TrimSpace(last),
		Alias:   strings.TrimSpace(last),
	}
}
