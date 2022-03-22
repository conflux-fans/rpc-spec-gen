package rust

import (
	"fmt"
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
)

func TestIsBaseType(t *testing.T) {
	r := IsBaseType("H256")
	assert.Equal(t, r, true)

	r = IsBaseType("crate::rpc::types::RpcAddress")
	assert.Equal(t, r, true)
}

func TestReadDir(t *testing.T) {
	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		fmt.Printf("f %s\n", f.Name())
	}
}
