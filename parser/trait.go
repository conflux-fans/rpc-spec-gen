package parser

import (
	"fmt"
	"regexp"
)

type RustTrait string

type RustTraitParsed struct {
	Comment string
	Name    string
	Funcs   []RustFuncParsed
}

func (rt RustTrait) Parse() RustTraitParsed {
	result := RustTraitParsed{}
	var traitReg = regexp.MustCompile(`(?mUs)(?:\/\/\/(.*\n)|)#\[rpc\(.*\)\]\npub trait (.*) \{(.*)}`)
	finds := traitReg.FindStringSubmatch(string(rt))

	result.Comment, result.Name = finds[1], finds[2]
	funcs := finds[3]

	fmt.Printf("traitRegFinded len %v, %v\n", len(finds), finds)
	fmt.Printf("RustTraitParsed %#v\n", result)
	fmt.Printf("funcs %v\n", funcs)

	var funcReg = regexp.MustCompile(`(?mUs)(?:\/\/\/.*\n|)\s*#\[rpc\(name =.*\)\].*;`)
	funcFinded := funcReg.FindAllString(funcs, -1)
	result.Funcs = make([]RustFuncParsed, len(funcFinded))
	for i, func_ := range funcFinded {
		fmt.Printf("func_ %v %v\n", i, func_)
		result.Funcs[i] = RustFunc(func_).Parse()
	}
	return result
}
