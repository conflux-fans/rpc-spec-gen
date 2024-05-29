package rust

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

type Trait string

type TraitParsed struct {
	Comment string
	Name    string
	Funcs   []FuncParsed
}

func (rt Trait) Parse() TraitParsed {
	result := TraitParsed{}
	var traitReg = regexp.MustCompile(`(?mUs)(?:\/\/\/(.*\n)|)#\[rpc\(.*\)\]\npub trait (.*) \{(.*)}`)
	finds := traitReg.FindStringSubmatch(string(rt))

	result.Comment, result.Name = finds[1], finds[2]
	funcs := finds[3]

	// fmt.Printf("traitRegFinded len %v, %v\n", len(finds), finds)
	// fmt.Printf("RustTraitParsed %#v\n", result)
	// fmt.Printf("funcs %v\n", funcs)

	logrus.WithField("trait regs num", len(finds)).WithField("rust trait parsed", result).WithField("funcs", funcs).Info("split trait to segments")

	var funcReg = regexp.MustCompile(`(?mUs)(?:\/\/\/.*\n|)\s*#\[rpc\(name =.*\)\].*;`)
	funcFinded := funcReg.FindAllString(funcs, -1)
	result.Funcs = make([]FuncParsed, len(funcFinded))
	for i, func_ := range funcFinded {
		fmt.Printf("func_ %v %v\n", i, func_)
		result.Funcs[i] = Func(func_).Parse()
	}
	return result
}
