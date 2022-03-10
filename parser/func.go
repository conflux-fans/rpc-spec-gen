package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type RustFunc string

type RustFuncParsed struct {
	Comment    string
	RpcMethod  string
	FuncName   string
	Params     []RustParamParsed
	ReturnType RustTypeParsed
}

type RustParamParsed struct {
	Name string
	Type RustTypeParsed
}

func (r RustFunc) Parse() RustFuncParsed {
	rfp := RustFuncParsed{}

	funcReg := regexp.MustCompile(`(?ims)(.*)?#\[rpc.*?"(.*)?"\)\].*?fn(.*)?\(\s*\&self,*(.*)?\)\s*->\s*.*?<(.*)>`)
	funcFinded := funcReg.FindStringSubmatch(string(r))

	// fmt.Printf("%#v\n", funcFinded[1:])
	comment, rpcMethod, funcName, params, returns := funcFinded[1], funcFinded[2], funcFinded[3], funcFinded[4], funcFinded[5]
	fmt.Printf("comment %v\nmethod %v\nfuncName %v\nparams %v\nreturns %v\n", comment, rpcMethod, funcName, params, returns)

	splitParams := strings.Split(params, ",")
	for _, param := range splitParams {
		fmt.Printf("param %v\n", param)
		if len(strings.TrimSpace(param)) == 0 {
			continue
		}

		paramReg := regexp.MustCompile(`(?ims)(.*): (.*)`)
		paramFinded := paramReg.FindStringSubmatch(param)
		fmt.Printf("paramFinded %v\n", paramFinded)

		name, type_ := paramFinded[1], paramFinded[2]
		fmt.Printf("name %v\ntype %v\n", name, type_)

		rfp.Params = append(rfp.Params, RustParamParsed{
			Name: name,
			Type: RustType(type_).Parse(),
		})
	}

	rfp.Comment = comment
	rfp.RpcMethod = rpcMethod
	rfp.FuncName = funcName
	rfp.ReturnType = RustType(returns).Parse()
	return rfp
}
