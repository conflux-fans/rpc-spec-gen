package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
)

func TestParseFunc(t *testing.T) {
	var rustFunc RustFunc = `/// Returns the number of uncles in a block with given block number.
     #[rpc(name = "eth_getUncleCountByBlockNumber")]
     fn block_uncles_count_by_number(
         &self, blockNum: BlockNumber,
     ) -> Result<Option<U256>>;`

	// // parsed := rustFunc.Parse()

	// // funcReg := regexp.MustCompile(`([\s\S]*)?#\[rpc.*?"(.*)?"\)\][\s\S]*?fn([\s\S]*)?\(\s*\&self,*([\s\S]*)?\)\s*->\s*.*?<(.*)>`)
	// funcReg := regexp.MustCompile(`(?ims)(.*)?#\[rpc.*?"(.*)?"\)\].*?fn(.*)?\(\s*\&self,*(.*)?\)\s*->\s*.*?<(.*)>`)
	// funcFinded := funcReg.FindStringSubmatch(string(rustFunc))

	// // fmt.Printf("%#v\n", funcFinded[1:])

	// comment, rpcMethod, funcName, params, returns := funcFinded[1], funcFinded[2], funcFinded[3], funcFinded[4], funcFinded[5]
	// fmt.Printf("comment %v\nmethod %v\nfuncName %v\nparams %v\nreturns %v\n", comment, rpcMethod, funcName, params, returns)

	// // params = strings.TrimSpace(params)
	// // params = strings.TrimSuffix(params, ",")
	// splitParams := strings.Split(params, ",")

	// for _, param := range splitParams {

	// 	fmt.Printf("param %v\n", param)
	// 	if len(strings.TrimSpace(param)) == 0 {
	// 		continue
	// 	}
	// 	paramReg := regexp.MustCompile(`(?ims)(.*): (.*)`)
	// 	paramFinded := paramReg.FindStringSubmatch(param)
	// 	fmt.Printf("paramFinded %v\n", paramFinded)

	// 	name, type_ := paramFinded[1], paramFinded[2]

	// 	fmt.Printf("name %v\ntype %v\n", name, type_)
	// }
	p := rustFunc.Parse()
	b, _ := json.MarshalIndent(p, "", "  ")
	fmt.Printf("%s\n", b)
}

func TestRegex(t *testing.T) {
	// reg := regexp.MustCompile(`(?U).*<([^>]*)>`)
	// finds := reg.FindStringSubmatch(`Result<Option<U256>>`)
	// fmt.Printf("%#v\n", finds)
	var re = regexp.MustCompile(`\s*.*?<(.*)>`)
	// var str = `Result<Option<U256>>`

	finds := re.FindStringSubmatch(`Result<Option<U256>>`)
	fmt.Printf("%#v\n", finds)

	// if len(re.FindStringIndex(str)) > 0 {
	// 	fmt.Println(re.FindString(str), "found at index", re.FindStringIndex(str)[0])
	// }
}
