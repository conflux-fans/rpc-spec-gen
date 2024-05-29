package rust

import (
	"regexp"
	"strings"

	"github.com/conflux-fans/rpc-spec-gen/utils"
	"github.com/sirupsen/logrus"
)

type Func string

type FuncParsed struct {
	Comment   string
	RpcMethod string
	FuncName  string
	Params    []ParamParsed
	Return    ReturnParsed
}

type ParamParsed struct {
	Name string
	Type TypeParsed
}

type ReturnParsed struct {
	Name string
	Type TypeParsed
}

func (r Func) Parse() FuncParsed {
	rfp := FuncParsed{}

	funcReg := regexp.MustCompile(`(?ims)(.*)?#\[rpc.*?"(.*)?"\)\].*?fn(.*)?\(\s*\&self,*(.*)?\)\s*->\s*.*?<(.*)>`)
	funcFinded := funcReg.FindStringSubmatch(string(r))

	// Note: returns 是不包含最外围泛型的结果，如RpcResult<H256>，这里只取H256
	comment, rpcMethod, funcName, params, returns := funcFinded[1], funcFinded[2], funcFinded[3], funcFinded[4], funcFinded[5]
	comment = utils.CleanComment(comment)
	// fmt.Printf("comment %v\nmethod %v\nfuncName %v\nparams %v\nreturns %v\n", comment, rpcMethod, funcName, params, returns)

	logrus.WithField("comment", comment).WithField("method", rpcMethod).WithField("func name", funcName).WithField("params", params).WithField("returns", returns).Info("split func to segments")

	splitParams := strings.Split(params, ",")
	for _, param := range splitParams {
		// fmt.Printf("param %v\n", param)
		if len(strings.TrimSpace(param)) == 0 {
			continue
		}

		paramReg := regexp.MustCompile(`(?ims)(.*): (.*)`)
		paramFinded := paramReg.FindStringSubmatch(param)
		// fmt.Printf("paramFinded %v\n", paramFinded)
		logrus.WithField("param", paramFinded).Debug("find param")

		name, type_ := paramFinded[1], paramFinded[2]
		// fmt.Printf("name %v\ntype %v\n", name, type_)
		logrus.WithField("name", name).WithField("type", type_).Debug("specific to name and type")

		rfp.Params = append(rfp.Params, ParamParsed{
			Name: strings.TrimSpace(name),
			Type: RustType(type_).Parse(),
		})
	}

	rfp.Comment = comment
	rfp.RpcMethod = rpcMethod
	rfp.FuncName = strings.TrimSpace(funcName)
	rfp.Return.Type = RustType(returns).Parse()
	return rfp
}
