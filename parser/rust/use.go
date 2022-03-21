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
}

func (u *Use) Parse() []UseType {
	body := u.body()
	nodes := body.toNodes()
	// fmt.Printf("nodes %#v\n", nodes)

	flattens := []string{}
	for _, node := range nodes {
		flattens = append(flattens, node.flatten()...)
	}

	uts := []UseType{}
	for i := range flattens {
		uts = append(uts, useItem(flattens[i]).Parse())
	}
	return uts
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
	str := strings.TrimSpace(string(*b))
	if str == "" {
		return nil
	}

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

func (ui useItem) Parse() UseType {
	nospaces := regexp.MustCompile(`\s+`).ReplaceAllString(string(ui), "")
	parts := strings.Split(nospaces, "::")
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

	return UseType{
		ModPath: parts[0 : len(parts)-1],
		Name:    strings.TrimSpace(last),
		Alias:   strings.TrimSpace(last),
	}
}
