package rust

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
)

func TestParseFunc(t *testing.T) {
	var rustFunc []Func = []Func{
		`/// Returns the number of uncles in a block with given block number.
		#[rpc(name = "eth_getUncleCountByBlockNumber")]
		fn block_uncles_count_by_number(
			&self, blockNum: BlockNumber,
		) -> Result<Option<U256>>;`,
	}

	p := rustFunc[0].Parse()
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
